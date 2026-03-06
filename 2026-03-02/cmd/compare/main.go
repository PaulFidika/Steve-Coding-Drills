package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Report struct {
	Backend  string      `json:"backend"`
	Scenario string      `json:"scenario"`
	Best     *StepResult `json:"best,omitempty"`
	RYOW     *RYOWResult `json:"ryw,omitempty"`
}

type StepResult struct {
	TargetRPS   int     `json:"target_rps"`
	AchievedRPS float64 `json:"achieved_rps"`
	P50MS       float64 `json:"p50_ms"`
	P95MS       float64 `json:"p95_ms"`
	P99MS       float64 `json:"p99_ms"`
	ErrorRate   float64 `json:"error_rate"`
}

type RYOWResult struct {
	Requests   int   `json:"requests"`
	Success    int64 `json:"success"`
	Failures   int64 `json:"failures"`
	Violations int64 `json:"violations"`
}

func main() {
	var (
		raftGet    = flag.String("raft-get", "results/raft-get.json", "raft GET report")
		raftPost   = flag.String("raft-post", "results/raft-post.json", "raft POST report")
		raftRYOW   = flag.String("raft-ryw", "results/raft-ryw.json", "raft RYOW report")
		garnetGet  = flag.String("garnet-get", "results/garnet-get.json", "garnet GET report")
		garnetPost = flag.String("garnet-post", "results/garnet-post.json", "garnet POST report")
		garnetRYOW = flag.String("garnet-ryw", "results/garnet-ryw.json", "garnet RYOW report")
		out        = flag.String("out", "results/comparison.md", "markdown output")
	)
	flag.Parse()

	rg := mustLoad(*raftGet)
	rp := mustLoad(*raftPost)
	rr := mustLoad(*raftRYOW)
	gg := mustLoad(*garnetGet)
	gp := mustLoad(*garnetPost)
	gr := mustLoad(*garnetRYOW)

	lines := []string{
		"# RAFT vs Garnet Benchmark",
		"",
		fmt.Sprintf("Generated: %s UTC", time.Now().UTC().Format(time.RFC3339)),
		"",
		"## Throughput and Latency (Best Sustainable Point)",
		"",
		"| Scenario | Backend | Target RPS | Achieved RPS | p50 (ms) | p95 (ms) | p99 (ms) | Error Rate |",
		"|---|---:|---:|---:|---:|---:|---:|---:|",
		row("GET", "raft", rg.Best),
		row("GET", "garnet", gg.Best),
		row("POST", "raft", rp.Best),
		row("POST", "garnet", gp.Best),
		"",
		"## Read-Your-Own-Write Consistency",
		"",
		"| Backend | Requests | Success | Failures | Violations |",
		"|---|---:|---:|---:|---:|",
		ryowRow("raft", rr.RYOW),
		ryowRow("garnet", gr.RYOW),
		"",
	}

	if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil {
		panic(err)
	}
	if err := os.WriteFile(*out, []byte(strings.Join(lines, "\n")), 0o644); err != nil {
		panic(err)
	}
	fmt.Printf("comparison written: %s\n", *out)
}

func mustLoad(path string) Report {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var r Report
	if err := json.Unmarshal(b, &r); err != nil {
		panic(err)
	}
	return r
}

func row(scenario, backend string, step *StepResult) string {
	if step == nil {
		return fmt.Sprintf("| %s | %s | n/a | n/a | n/a | n/a | n/a | n/a |", scenario, backend)
	}
	return fmt.Sprintf("| %s | %s | %d | %.1f | %.2f | %.2f | %.2f | %.4f |", scenario, backend, step.TargetRPS, step.AchievedRPS, step.P50MS, step.P95MS, step.P99MS, step.ErrorRate)
}

func ryowRow(backend string, ryow *RYOWResult) string {
	if ryow == nil {
		return fmt.Sprintf("| %s | n/a | n/a | n/a | n/a |", backend)
	}
	return fmt.Sprintf("| %s | %d | %d | %d | %d |", backend, ryow.Requests, ryow.Success, ryow.Failures, ryow.Violations)
}
