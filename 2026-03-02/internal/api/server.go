package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"raft-garnet-bench/internal/backend"
	"raft-garnet-bench/internal/netutil"
)

type Server struct {
	listenAddr string
	nodeID     uint64
	store      backend.Store
	client     *http.Client
	httpServer *http.Server
}

type writeRequest struct {
	Value string `json:"value"`
}

type stateResponse struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	NodeID uint64 `json:"node_id"`
	Mode   string `json:"mode"`
}

const (
	forwardCountHeader = "X-Raft-Forwarded-Count"
	maxForwardHops     = 3
)

func New(listenAddr string, nodeID uint64, store backend.Store) *Server {
	return &Server{
		listenAddr: listenAddr,
		nodeID:     nodeID,
		store:      store,
		client: netutil.NewHTTPClient(netutil.HTTPClientOptions{
			Timeout:             5 * time.Second,
			MaxIdleConns:        1024,
			MaxIdleConnsPerHost: 512,
			MaxConnsPerHost:     512,
		}),
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/state/", s.handleState)

	s.httpServer = &http.Server{
		Addr:              s.listenAddr,
		Handler:           mux,
		ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	log.Printf("[server] listening on %s node_id=%d mode=%s", s.listenAddr, s.nodeID, s.store.Mode())
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":      true,
		"node_id": s.nodeID,
		"mode":    s.store.Mode(),
	})
}

func (s *Server) handleState(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimPrefix(r.URL.Path, "/state/")
	if key == "" || strings.Contains(key, "/") {
		writeError(w, http.StatusBadRequest, "invalid key")
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.getState(w, r, key)
	case http.MethodPost:
		s.postState(w, r, key)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Server) getState(w http.ResponseWriter, r *http.Request, key string) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	value, err := s.store.Get(ctx, key)
	if err == nil {
		writeJSON(w, http.StatusOK, stateResponse{Key: key, Value: value, NodeID: s.nodeID, Mode: s.store.Mode()})
		return
	}

	if errors.Is(err, backend.ErrNotFound) {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	if nle, ok := err.(*backend.NotLeaderError); ok {
		if !s.canForward(r) {
			writeError(w, http.StatusServiceUnavailable, "leader changed during forward")
			return
		}
		if err := s.forwardToLeader(w, r, nle); err != nil {
			writeError(w, http.StatusServiceUnavailable, err.Error())
		}
		return
	}
	writeError(w, http.StatusInternalServerError, err.Error())
}

func (s *Server) postState(w http.ResponseWriter, r *http.Request, key string) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed reading body")
		return
	}
	defer r.Body.Close()

	value := parseValue(body)
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := s.store.Set(ctx, key, value); err == nil {
		writeJSON(w, http.StatusOK, stateResponse{Key: key, Value: value, NodeID: s.nodeID, Mode: s.store.Mode()})
		return
	} else {
		if nle, ok := err.(*backend.NotLeaderError); ok {
			if !s.canForward(r) {
				writeError(w, http.StatusServiceUnavailable, "leader changed during forward")
				return
			}
			if err := s.forwardToLeader(w, rebuildRequestBody(r, body), nle); err != nil {
				writeError(w, http.StatusServiceUnavailable, err.Error())
			}
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
	}
}

func (s *Server) canForward(r *http.Request) bool {
	count := 0
	if raw := strings.TrimSpace(r.Header.Get(forwardCountHeader)); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 {
			count = n
		}
	}
	return count < maxForwardHops
}

func (s *Server) forwardToLeader(w http.ResponseWriter, r *http.Request, nle *backend.NotLeaderError) error {
	if nle == nil || strings.TrimSpace(nle.LeaderURL) == "" {
		return fmt.Errorf("leader unavailable")
	}
	target := strings.TrimSuffix(nle.LeaderURL, "/") + r.URL.Path
	if r.URL.RawQuery != "" {
		target += "?" + r.URL.RawQuery
	}

	var body io.Reader
	if r.Body != nil {
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("read body for forwarding: %w", err)
		}
		r.Body = io.NopCloser(bytes.NewReader(payload))
		body = bytes.NewReader(payload)
	}

	fr, err := http.NewRequestWithContext(r.Context(), r.Method, target, body)
	if err != nil {
		return fmt.Errorf("build forward request: %w", err)
	}
	fr.Header = r.Header.Clone()
	forwardCount := 0
	if raw := strings.TrimSpace(r.Header.Get(forwardCountHeader)); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n > 0 {
			forwardCount = n
		}
	}
	fr.Header.Set(forwardCountHeader, strconv.Itoa(forwardCount+1))

	resp, err := s.client.Do(fr)
	if err != nil {
		return fmt.Errorf("forward request failed: %w", err)
	}
	defer resp.Body.Close()

	for k, vals := range resp.Header {
		for _, v := range vals {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
	return nil
}

func parseValue(body []byte) string {
	if len(body) == 0 {
		return ""
	}
	var req writeRequest
	if err := json.Unmarshal(body, &req); err == nil {
		return req.Value
	}
	return string(body)
}

func rebuildRequestBody(r *http.Request, body []byte) *http.Request {
	r2 := r.Clone(r.Context())
	r2.Body = io.NopCloser(bytes.NewReader(body))
	return r2
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]any{"error": msg})
}

func writeJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}
