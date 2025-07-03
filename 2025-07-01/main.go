package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

// static list of objects present in both buckets
var objectKeys = []string{
	"t-myyt993e.jpg",
	"t-66ns5cax.jpg",
	"n-myyt993e.jpg",
	"n-66ns5cax.jpg",
	"n-3bxtmaxs.jpg",
	"t-3bxtmaxs.jpg",
}

// JSON response structure
type imageList struct {
	URLs []string `json:"urls"`
}

func main() {
	// load .env file if present
	_ = godotenv.Load()

	mode := flag.String("mode", "server", "mode: server | bench")
	addr := flag.String("addr", ":8080", "server listen address (server mode)")
	serverBase := flag.String("server", "http://localhost:8080", "server base URL (bench mode)")
	runs := flag.Int("runs", 3, "number of benchmark repetitions (bench mode)")
	flag.Parse()

	switch *mode {
	case "server":
		startServer(*addr)
	case "bench":
		runBench(*serverBase, *runs)
	default:
		log.Fatalf("unknown mode %s", *mode)
	}
}

// ---------------------- SERVER ----------------------
func startServer(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/images/public", publicHandler)
	mux.HandleFunc("/images/public_raw", publicRawHandler)
	mux.HandleFunc("/images/private", privateHandler)

	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

// returns public URLs under https://test1.doujins.media
func publicHandler(w http.ResponseWriter, r *http.Request) {
	base := "https://test1.doujins.media/"
	// add per-response nonce to prevent client/proxy caching of the objects
	nonce := fmt.Sprintf("%d", time.Now().UnixNano())
	urls := make([]string, 0, len(objectKeys))
	for _, k := range objectKeys {
		urls = append(urls, base+k+"?n="+nonce)
	}
	respondJSON(w, imageList{URLs: urls})
}

// returns raw R2 URLs for public bucket (no Cloudflare CDN)
func publicRawHandler(w http.ResponseWriter, r *http.Request) {
	base := os.Getenv("PUBLIC_RAW_BASE")
	if base == "" {
		base = "https://b83d3b9a7092a7da1bee02c72968b81e.r2.cloudflarestorage.com/test1/"
	}
	if !strings.HasSuffix(base, "/") {
		base += "/"
	}
	// cache-buster nonce (same per response to keep list size unchanged)
	nonce := fmt.Sprintf("%d", time.Now().UnixNano())
	urls := make([]string, 0, len(objectKeys))
	for _, k := range objectKeys {
		urls = append(urls, base+k+"?n="+nonce)
	}
	respondJSON(w, imageList{URLs: urls})
}

// returns presigned private URLs for R2 bucket
func privateHandler(w http.ResponseWriter, r *http.Request) {
	signer, err := newPresigner()
	if err != nil {
		http.Error(w, "signer setup failed", http.StatusInternalServerError)
		log.Printf("signer setup failed: %v", err)
		return
	}

	bucket := os.Getenv("R2_BUCKET")
	if bucket == "" {
		bucket = "test2" // fallback to default from prompt
	}

	ctx := r.Context()
	urls := make([]string, 0, len(objectKeys))
	for _, k := range objectKeys {
		req, _ := signer.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &k,
		}, s3.WithPresignExpires(15*time.Minute))
		urls = append(urls, req.URL)
	}
	respondJSON(w, imageList{URLs: urls})
}

// helper to write JSON with sane headers
func respondJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if err := enc.Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ---------------------- BENCH ----------------------
func runBench(serverBase string, runs int) {
	pubEndpoint := strings.TrimRight(serverBase, "/") + "/images/public"
	pubRawEndpoint := strings.TrimRight(serverBase, "/") + "/images/public_raw"
	privEndpoint := strings.TrimRight(serverBase, "/") + "/images/private"

	for i := 1; i <= runs; i++ {
		// measure raw first, then CDN
		rawDur := measure(pubRawEndpoint)
		cdnDur := measure(pubEndpoint)
		privDur := measure(privEndpoint)
		fmt.Printf("Run %d\tcdn: %v\traw: %v (Δ=%v)\tprivate: %v (Δ=%v)\n",
			i, cdnDur, rawDur, rawDur-cdnDur, privDur, privDur-cdnDur)
	}
}

func measure(listEndpoint string) time.Duration {
	start := time.Now()

	// fetch list
	resp, err := http.Get(listEndpoint)
	if err != nil {
		log.Fatalf("failed to fetch %s: %v", listEndpoint, err)
	}
	defer resp.Body.Close()
	var lst imageList
	if err := json.NewDecoder(resp.Body).Decode(&lst); err != nil {
		log.Fatalf("decode list: %v", err)
	}

	// fetch each resource concurrently
	var wg sync.WaitGroup
	errCh := make(chan error, len(lst.URLs))

	for _, u := range lst.URLs {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			resp, err := http.Get(url)
			if err != nil {
				errCh <- fmt.Errorf("fetch %s: %v", url, err)
				return
			}
			_, _ = io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}(u)
	}

	wg.Wait()
	close(errCh)
	if err, ok := <-errCh; ok {
		log.Fatalf("%v", err)
	}

	return time.Since(start)
}

// ---------------------- R2 signing ----------------------
// newPresigner sets up an AWS S3 V4 presigner for Cloudflare R2 using ENV vars:
// R2_ACCESS_KEY_ID, R2_SECRET_ACCESS_KEY, R2_ENDPOINT, R2_REGION (optional), R2_BUCKET (optional)
func newPresigner() (*s3.PresignClient, error) {
	accessKey := os.Getenv("R2_ACCESS_KEY_ID")
	secretKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	if accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("missing R2 credentials in env")
	}

	endpoint := os.Getenv("R2_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://b83d3b9a7092a7da1bee02c72968b81e.r2.cloudflarestorage.com" // default from prompt
	}

	region := os.Getenv("R2_REGION")
	if region == "" {
		region = "auto" // Cloudflare recommends the magic region "auto"
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, r string, _ ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:               endpoint,
				SigningRegion:     region,
				HostnameImmutable: true,
			}, nil
		})),
	)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true // R2 requires path-style addressing
	})
	return s3.NewPresignClient(s3Client), nil
}
