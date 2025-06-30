package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/redis/go-redis/v9"
)

const totalKeys = 200 // change this to adjust how many objects to sign

func main() {
	ctx := context.Background()

	bucket := os.Getenv("BUCKET_NAME")
	if bucket == "" {
		bucket = "monkey"
		// log.Fatal("BUCKET_NAME environment variable must be set")
	}

	rand.Seed(time.Now().UnixNano())

	keys := make([]string, totalKeys)
	fixedCount := totalKeys / 2 // 50% deterministic / 50% random
	for i := 0; i < fixedCount; i++ {
		keys[i] = fmt.Sprintf("image_fixed_%03d.jpg", i)
	}
	for i := fixedCount; i < totalKeys; i++ {
		keys[i] = fmt.Sprintf("image_rand_%06d.jpg", rand.Intn(1_000_000))
	}

	presigner, err := newPresignClient(ctx)
	if err != nil {
		log.Fatalf("unable to create presign client: %v", err)
	}

	// --- Redis client (assumes local Redis on default port) ---
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6380",
		DB:   0,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("Redis unavailable: %v (skipping cache benchmark)", err)
	}

	start := time.Now()
	if _, err := signSequential(ctx, presigner, bucket, keys); err != nil {
		log.Fatalf("sequential signing failed: %v", err)
	}
	seqDuration := time.Since(start)

	const workers = 10

	start = time.Now()
	if _, err := signConcurrentChunks(ctx, presigner, bucket, keys, workers); err != nil {
		log.Fatalf("concurrent signing failed: %v", err)
	}
	concDuration := time.Since(start)

	fmt.Printf("Sequential signing took: %v\n", seqDuration)
	fmt.Printf("Concurrent signing took: %v\n", concDuration)

	if rdb != nil {
		start = time.Now()
		if _, err := cacheOrSignBatch(ctx, rdb, presigner, bucket, keys); err != nil {
			log.Fatalf("redis signing failed: %v", err)
		}
		cacheDuration := time.Since(start)
		fmt.Printf("Redis cache retrieval/signing took: %v\n", cacheDuration)
	}
}

func newPresignClient(ctx context.Context) (*s3.PresignClient, error) {
	endpoint := os.Getenv("S3_ENDPOINT")
	region := os.Getenv("AWS_REGION")
	if region == "" {
		// Cloudflare R2 typically uses "auto" for the signing region.
		region = "auto"
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	// Override endpoint when provided (useful for Cloudflare R2).
	if endpoint != "" {
		cfg.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(func(service, r string, opts ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:           endpoint,
				SigningRegion: region,
			}, nil
		})
	}

	client := s3.NewFromConfig(cfg)
	return s3.NewPresignClient(client), nil
}

func signSequential(ctx context.Context, presigner *s3.PresignClient, bucket string, keys []string) ([]string, error) {
	urls := make([]string, len(keys))
	for i, key := range keys {
		out, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &key,
		}, s3.WithPresignExpires(time.Hour))
		if err != nil {
			return nil, err
		}
		urls[i] = out.URL
	}
	return urls, nil
}

func signConcurrentChunks(ctx context.Context, presigner *s3.PresignClient, bucket string, keys []string, workers int) ([]string, error) {
	if workers < 1 {
		workers = 1
	}

	type piece struct {
		start int
		urls  []string
		err   error
	}

	n := len(keys)
	urls := make([]string, n)

	// Determine chunk boundaries (ceil division)
	chunk := (n + workers - 1) / workers

	piecesCh := make(chan piece, workers)

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		start := w * chunk
		if start >= n {
			break
		}
		end := start + chunk
		if end > n {
			end = n
		}

		wg.Add(1)
		go func(s, e int) {
			defer wg.Done()
			local := make([]string, e-s)
			for i := s; i < e; i++ {
				key := keys[i]
				out, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
					Bucket: &bucket,
					Key:    &key,
				}, s3.WithPresignExpires(time.Hour))
				if err != nil {
					piecesCh <- piece{err: err}
					return
				}
				local[i-s] = out.URL
			}
			piecesCh <- piece{start: s, urls: local}
		}(start, end)
	}

	// Close channel once all goroutines have finished
	go func() {
		wg.Wait()
		close(piecesCh)
	}()

	for p := range piecesCh {
		if p.err != nil {
			return nil, p.err
		}
		copy(urls[p.start:], p.urls)
	}

	return urls, nil
}

// cacheOrSignBatch fetches all URLs with a single MGET. Any missing entries are
// signed, then cached with a pipelined MSET (one round-trip). This avoids the
// per-key latency overhead that hurt the previous benchmark.
func cacheOrSignBatch(ctx context.Context, rdb *redis.Client, presigner *s3.PresignClient, bucket string, keys []string) ([]string, error) {
	n := len(keys)
	urls := make([]string, n)

	// 1. Bulk fetch
	vals, err := rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	// 2. Track misses
	type miss struct{ idx int; key string }
	var misses []miss
	for i, v := range vals {
		if v == nil {
			misses = append(misses, miss{idx: i, key: keys[i]})
			continue
		}
		urls[i] = v.(string)
	}

	// 3. Sign misses concurrently to reduce CPU time.
	if len(misses) > 0 {
		// Parallel sign
		const signWorkers = 10
		signed := make([]string, len(misses))
		var idx64 int64
		var wg sync.WaitGroup
		var signErr error

		for w := 0; w < signWorkers; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					i := int(atomic.AddInt64(&idx64, 1)) - 1
					if i >= len(misses) {
						return
					}
					m := misses[i]
					out, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
						Bucket: &bucket,
						Key:    &m.key,
					}, s3.WithPresignExpires(time.Hour))
					if err != nil {
						signErr = err
						return
					}
					signed[i] = out.URL
				}
			}()
		}
		wg.Wait()
		if signErr != nil {
			return nil, signErr
		}

		// 4. Queue SETEX commands in a pipeline (no Lua).
		// Build args for MSET (key1 val1 key2 val2 ...)
		args := make([]interface{}, 0, len(misses)*2)
		for i, m := range misses {
			url := signed[i]
			urls[m.idx] = url
			args = append(args, m.key, url)
		}

		pipe := rdb.Pipeline()
		pipe.MSet(ctx, args...)
		const ttl = time.Hour
		for _, m := range misses {
			pipe.PExpire(ctx, m.key, ttl)
		}
		if _, err := pipe.Exec(ctx); err != nil {
			return nil, err
		}
	}

	return urls, nil
}