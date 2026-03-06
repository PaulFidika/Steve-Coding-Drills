package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"raft-garnet-bench/internal/netutil"
)

type Report struct {
	Backend     string       `json:"backend"`
	BaseURL     string       `json:"base_url"`
	Scenario    string       `json:"scenario"`
	GeneratedAt time.Time    `json:"generated_at"`
	Config      BenchConfig  `json:"config"`
	Steps       []StepResult `json:"steps,omitempty"`
	Best        *StepResult  `json:"best,omitempty"`
	RYOW        *RYOWResult  `json:"ryw,omitempty"`
}

type BenchConfig struct {
	DurationSeconds float64 `json:"duration_seconds"`
	Workers         int     `json:"workers"`
	StartRPS        int     `json:"start_rps"`
	StepRPS         int     `json:"step_rps"`
	MaxRPS          int     `json:"max_rps"`
	Keyspace        int     `json:"keyspace"`
	SeedKeys        int     `json:"seed_keys"`
	TimeoutMS       int64   `json:"timeout_ms"`
	ErrorThreshold  float64 `json:"error_threshold"`
	P95ThresholdMS  float64 `json:"p95_threshold_ms"`
	RoutingMode     string  `json:"routing_mode"`
	NodeURLCount    int     `json:"node_url_count"`
	ShardCount      int     `json:"shard_count"`
	ReplicationFact int     `json:"replication_factor"`
}

type StepResult struct {
	TargetRPS   int            `json:"target_rps"`
	AchievedRPS float64        `json:"achieved_rps"`
	Sent        int64          `json:"sent"`
	Dropped     int64          `json:"dropped"`
	Success     int64          `json:"success"`
	Failures    int64          `json:"failures"`
	ErrorRate   float64        `json:"error_rate"`
	P50MS       float64        `json:"p50_ms"`
	P95MS       float64        `json:"p95_ms"`
	P99MS       float64        `json:"p99_ms"`
	StatusCodes map[string]int `json:"status_codes"`
}

type RYOWResult struct {
	Requests        int     `json:"requests"`
	Concurrency     int     `json:"concurrency"`
	Success         int64   `json:"success"`
	Failures        int64   `json:"failures"`
	Violations      int64   `json:"violations"`
	DurationSeconds float64 `json:"duration_seconds"`
}

type targetRouter struct {
	mode              string
	baseURL           string
	nodeURLs          []string
	shardCount        int
	replicationFactor int
}

func main() {
	var (
		baseURL        = flag.String("base-url", "http://localhost:18080", "base URL of cluster LB")
		nodeURLsFlag   = flag.String("node-urls", "", "comma-separated node URLs used by random/shard routing modes")
		routingMode    = flag.String("routing", "single", "routing mode: single|random|shard")
		routeShards    = flag.Int("route-shard-count", 3, "shard count used by shard-aware router")
		routeRF        = flag.Int("route-replication-factor", 3, "replication factor used by shard-aware router")
		backend        = flag.String("backend", "unknown", "backend label for report")
		scenario       = flag.String("scenario", "post", "one of: get, post, mix, ryow")
		duration       = flag.Duration("duration", 20*time.Second, "duration per throughput step")
		workers        = flag.Int("workers", 128, "worker goroutines")
		startRPS       = flag.Int("start-rps", 200, "initial target requests/sec")
		stepRPS        = flag.Int("step-rps", 200, "increment requests/sec each step")
		maxRPS         = flag.Int("max-rps", 5000, "max target requests/sec")
		keyspace       = flag.Int("keyspace", 10000, "number of keys used by load generator")
		seedKeys       = flag.Int("seed-keys", 10000, "keys to pre-seed for GET workloads")
		timeout        = flag.Duration("timeout", 3*time.Second, "http timeout per request")
		reportPath     = flag.String("report", "", "optional output report path (json)")
		errThreshold   = flag.Float64("error-threshold", 0.01, "max acceptable error rate for sustainable throughput")
		p95ThresholdMS = flag.Float64("p95-threshold-ms", 250, "max acceptable p95 latency for sustainable throughput")
		ryowRequests   = flag.Int("ryw-requests", 2000, "RYOW requests when scenario=ryow")
		rywConc        = flag.Int("ryw-concurrency", 64, "RYOW concurrency when scenario=ryow")
	)
	flag.Parse()

	cleanBase := strings.TrimRight(*baseURL, "/")
	router, err := newTargetRouter(*routingMode, cleanBase, *nodeURLsFlag, *routeShards, *routeRF)
	if err != nil {
		log.Fatalf("router config error: %v", err)
	}
	client := netutil.NewHTTPClient(netutil.HTTPClientOptions{
		Timeout:             *timeout,
		MaxIdleConns:        *workers * 16,
		MaxIdleConnsPerHost: *workers * 8,
		MaxConnsPerHost:     *workers * 16,
	})

	report := Report{
		Backend:     *backend,
		BaseURL:     router.baseURL,
		Scenario:    *scenario,
		GeneratedAt: time.Now().UTC(),
		Config: BenchConfig{
			DurationSeconds: duration.Seconds(),
			Workers:         *workers,
			StartRPS:        *startRPS,
			StepRPS:         *stepRPS,
			MaxRPS:          *maxRPS,
			Keyspace:        *keyspace,
			SeedKeys:        *seedKeys,
			TimeoutMS:       timeout.Milliseconds(),
			ErrorThreshold:  *errThreshold,
			P95ThresholdMS:  *p95ThresholdMS,
			RoutingMode:     router.mode,
			NodeURLCount:    len(router.nodeURLs),
			ShardCount:      router.shardCount,
			ReplicationFact: router.replicationFactor,
		},
	}

	switch strings.ToLower(*scenario) {
	case "ryw":
		ryow, err := runRYOW(client, router, *ryowRequests, *rywConc)
		if err != nil {
			log.Fatalf("ryw failed: %v", err)
		}
		report.RYOW = &ryow
		fmt.Printf("RYOW requests=%d success=%d failures=%d violations=%d duration=%.2fs\n", ryow.Requests, ryow.Success, ryow.Failures, ryow.Violations, ryow.DurationSeconds)
	case "get", "post", "mix":
		if *scenario == "get" || *scenario == "mix" {
			seedCount := *seedKeys
			if seedCount <= 0 {
				seedCount = *keyspace
			}
			log.Printf("seeding %d keys", seedCount)
			if err := seed(client, router, seedCount); err != nil {
				log.Fatalf("seed failed: %v", err)
			}
		}

		for rps := *startRPS; rps <= *maxRPS; rps += *stepRPS {
			step := runStep(client, router, strings.ToLower(*scenario), *duration, rps, *workers, *keyspace)
			report.Steps = append(report.Steps, step)
			fmt.Printf("step rps=%d achieved=%.1f success=%d failures=%d err_rate=%.4f p95=%.2fms\n", step.TargetRPS, step.AchievedRPS, step.Success, step.Failures, step.ErrorRate, step.P95MS)

			if step.ErrorRate <= *errThreshold && step.P95MS <= *p95ThresholdMS {
				cp := step
				report.Best = &cp
			}
		}

		if report.Best != nil {
			fmt.Printf("best sustainable throughput: rps=%d achieved=%.1f p95=%.2fms error_rate=%.4f\n", report.Best.TargetRPS, report.Best.AchievedRPS, report.Best.P95MS, report.Best.ErrorRate)
		} else {
			fmt.Println("no sustainable throughput point met thresholds")
		}
	default:
		log.Fatalf("unsupported scenario %q", *scenario)
	}

	if *reportPath != "" {
		if err := writeReport(*reportPath, report); err != nil {
			log.Fatalf("write report: %v", err)
		}
		fmt.Printf("report written: %s\n", *reportPath)
	}
}

func newTargetRouter(mode string, baseURL string, nodeURLsCSV string, shardCount int, rf int) (*targetRouter, error) {
	r := &targetRouter{
		mode:              strings.ToLower(strings.TrimSpace(mode)),
		baseURL:           strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		shardCount:        shardCount,
		replicationFactor: rf,
	}
	if r.mode == "" {
		r.mode = "single"
	}
	if nodeURLsCSV != "" {
		parts := strings.Split(nodeURLsCSV, ",")
		r.nodeURLs = make([]string, 0, len(parts))
		for _, p := range parts {
			u := strings.TrimRight(strings.TrimSpace(p), "/")
			if u != "" {
				r.nodeURLs = append(r.nodeURLs, u)
			}
		}
	}
	switch r.mode {
	case "single":
		if r.baseURL == "" {
			return nil, fmt.Errorf("base-url is required for single routing")
		}
	case "random":
		if len(r.nodeURLs) == 0 {
			return nil, fmt.Errorf("node-urls is required for random routing")
		}
	case "shard":
		if len(r.nodeURLs) == 0 {
			return nil, fmt.Errorf("node-urls is required for shard routing")
		}
		if r.shardCount <= 0 {
			return nil, fmt.Errorf("route-shard-count must be > 0")
		}
		if r.replicationFactor <= 0 || r.replicationFactor > len(r.nodeURLs) {
			return nil, fmt.Errorf("route-replication-factor must be in 1..%d", len(r.nodeURLs))
		}
	default:
		return nil, fmt.Errorf("unsupported routing mode %q", r.mode)
	}
	if r.baseURL == "" && len(r.nodeURLs) > 0 {
		r.baseURL = r.nodeURLs[0]
	}
	return r, nil
}

func (r *targetRouter) urlForKey(key string, rnd *rand.Rand) string {
	switch r.mode {
	case "single":
		return r.baseURL
	case "random":
		if rnd == nil {
			rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
		}
		return r.nodeURLs[rnd.Intn(len(r.nodeURLs))]
	case "shard":
		shardID := int(computeShardID(key, uint64(r.shardCount)))
		replicas := shardReplicas(shardID, len(r.nodeURLs), r.shardCount, r.replicationFactor)
		if len(replicas) == 0 {
			return r.nodeURLs[0]
		}
		leaderID := replicas[0]
		if leaderID < 1 || leaderID > len(r.nodeURLs) {
			return r.nodeURLs[0]
		}
		return r.nodeURLs[leaderID-1]
	default:
		return r.baseURL
	}
}

func shardReplicas(shardID, nodes, shardCount, rf int) []int {
	if nodes <= 0 || shardCount <= 0 || rf <= 0 {
		return nil
	}
	if shardID < 1 {
		shardID = 1
	}
	leaderPos := ((shardID - 1) * nodes) / shardCount
	out := make([]int, 0, rf)
	for offset := 0; offset < nodes && len(out) < rf; offset++ {
		nodeID := ((leaderPos + offset) % nodes) + 1
		out = append(out, nodeID)
	}
	return out
}

func seed(client *http.Client, router *targetRouter, n int) error {
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("seed-%d", i)
		if err := postValue(client, router, key, value, nil); err != nil {
			return err
		}
	}
	return nil
}

func runStep(client *http.Client, router *targetRouter, scenario string, duration time.Duration, rps, workers, keyspace int) StepResult {
	if workers <= 0 {
		workers = 1
	}
	if keyspace <= 0 {
		keyspace = 1
	}

	jobs := make(chan int64, workers*4)
	start := time.Now()
	end := start.Add(duration)

	var metricMu sync.Mutex
	latencies := make([]int64, 0, rps*int(math.Ceil(duration.Seconds())))
	statusCodes := map[string]int{}
	var success int64
	var failures int64

	record := func(latency time.Duration, status int, failed bool) {
		metricMu.Lock()
		latencies = append(latencies, latency.Nanoseconds())
		statusCodes[fmt.Sprintf("%d", status)]++
		metricMu.Unlock()
		if failed {
			atomic.AddInt64(&failures, 1)
		} else {
			atomic.AddInt64(&success, 1)
		}
	}

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			rnd := rand.New(rand.NewSource(time.Now().UnixNano() + int64(workerID)*7919))
			for seq := range jobs {
				method, key, value := nextRequest(scenario, rnd, seq, keyspace)
				started := time.Now()
				status, err := execute(client, router, method, key, value, rnd)
				record(time.Since(started), status, err != nil || status < 200 || status >= 300)
			}
		}(i)
	}

	var sent int64
	var dropped int64
	var generated int64
	for {
		now := time.Now()
		if now.After(end) {
			break
		}
		expected := int64(now.Sub(start).Seconds() * float64(rps))
		for generated < expected {
			generated++
			select {
			case jobs <- generated:
				atomic.AddInt64(&sent, 1)
			default:
				atomic.AddInt64(&dropped, 1)
			}
		}
		time.Sleep(1 * time.Millisecond)
	}
	close(jobs)
	wg.Wait()

	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })
	successCount := atomic.LoadInt64(&success)
	failureCount := atomic.LoadInt64(&failures)
	sentCount := atomic.LoadInt64(&sent)
	droppedCount := atomic.LoadInt64(&dropped)
	total := successCount + failureCount
	errRate := 0.0
	if total > 0 {
		errRate = float64(failureCount) / float64(total)
	}

	elapsed := time.Since(start).Seconds()
	if elapsed <= 0 {
		elapsed = duration.Seconds()
	}

	return StepResult{
		TargetRPS:   rps,
		AchievedRPS: float64(successCount) / elapsed,
		Sent:        sentCount,
		Dropped:     droppedCount,
		Success:     successCount,
		Failures:    failureCount,
		ErrorRate:   errRate,
		P50MS:       percentileMS(latencies, 0.50),
		P95MS:       percentileMS(latencies, 0.95),
		P99MS:       percentileMS(latencies, 0.99),
		StatusCodes: statusCodes,
	}
}

func nextRequest(scenario string, rnd *rand.Rand, seq int64, keyspace int) (method, key, value string) {
	switch scenario {
	case "get":
		key = fmt.Sprintf("key-%d", rnd.Intn(keyspace))
		return http.MethodGet, key, ""
	case "mix":
		if rnd.Intn(2) == 0 {
			key = fmt.Sprintf("key-%d", rnd.Intn(keyspace))
			return http.MethodGet, key, ""
		}
		fallthrough
	case "post":
		key = fmt.Sprintf("key-%d", rnd.Intn(keyspace))
		value = fmt.Sprintf("v-%d-%d", seq, rnd.Int63())
		return http.MethodPost, key, value
	default:
		key = fmt.Sprintf("key-%d", rnd.Intn(keyspace))
		value = fmt.Sprintf("v-%d-%d", seq, rnd.Int63())
		return http.MethodPost, key, value
	}
}

func execute(client *http.Client, router *targetRouter, method, key, value string, rnd *rand.Rand) (int, error) {
	url := fmt.Sprintf("%s/state/%s", router.urlForKey(key, rnd), key)
	if method == http.MethodGet {
		resp, err := client.Get(url)
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
		return resp.StatusCode, nil
	}

	body, _ := json.Marshal(map[string]string{"value": value})
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	return resp.StatusCode, nil
}

func runRYOW(client *http.Client, router *targetRouter, requests, concurrency int) (RYOWResult, error) {
	if requests <= 0 {
		requests = 1
	}
	if concurrency <= 0 {
		concurrency = 1
	}
	jobs := make(chan int, concurrency*2)
	var success int64
	var failures int64
	var violations int64

	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			rnd := rand.New(rand.NewSource(time.Now().UnixNano() + int64(workerID)*31337))
			for n := range jobs {
				key := fmt.Sprintf("ryw-%d-%d-%d", workerID, n, rnd.Int63())
				value := fmt.Sprintf("value-%d", rnd.Int63())
				if err := postValue(client, router, key, value, rnd); err != nil {
					atomic.AddInt64(&failures, 1)
					continue
				}
				got, err := getValue(client, router, key, rnd)
				if err != nil {
					atomic.AddInt64(&failures, 1)
					continue
				}
				if got != value {
					atomic.AddInt64(&violations, 1)
					continue
				}
				atomic.AddInt64(&success, 1)
			}
		}(i)
	}

	for i := 0; i < requests; i++ {
		jobs <- i
	}
	close(jobs)
	wg.Wait()

	return RYOWResult{
		Requests:        requests,
		Concurrency:     concurrency,
		Success:         success,
		Failures:        failures,
		Violations:      violations,
		DurationSeconds: time.Since(start).Seconds(),
	}, nil
}

func postValue(client *http.Client, router *targetRouter, key, value string, rnd *rand.Rand) error {
	url := fmt.Sprintf("%s/state/%s", router.urlForKey(key, rnd), key)
	body, _ := json.Marshal(map[string]string{"value": value})
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("post status=%d", resp.StatusCode)
	}
	return nil
}

func getValue(client *http.Client, router *targetRouter, key string, rnd *rand.Rand) (string, error) {
	url := fmt.Sprintf("%s/state/%s", router.urlForKey(key, rnd), key)
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("get status=%d", resp.StatusCode)
	}
	var parsed struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", err
	}
	return parsed.Value, nil
}

func computeShardID(key string, shardCount uint64) uint64 {
	if shardCount == 0 {
		return 1
	}
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return uint64(h.Sum32()%uint32(shardCount)) + 1
}

func percentileMS(values []int64, p float64) float64 {
	if len(values) == 0 {
		return 0
	}
	idx := int(math.Ceil(float64(len(values))*p)) - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(values) {
		idx = len(values) - 1
	}
	return float64(values[idx]) / float64(time.Millisecond)
}

func writeReport(path string, report Report) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
