package raftstore

import (
	"encoding/json"
	"fmt"
	"io"

	sm "github.com/lni/dragonboat/v4/statemachine"
)

type stateMachine struct {
	kv map[string]string
}

func newStateMachine() sm.IStateMachine {
	return &stateMachine{kv: map[string]string{}}
}

func (s *stateMachine) Update(entry sm.Entry) (sm.Result, error) {
	var c command
	if err := json.Unmarshal(entry.Cmd, &c); err != nil {
		return sm.Result{}, nil
	}
	switch c.Type {
	case "set":
		var payload setCommand
		if err := json.Unmarshal(c.Payload, &payload); err == nil && payload.Key != "" {
			s.kv[payload.Key] = payload.Value
		}
	}
	return sm.Result{Value: 1}, nil
}

func (s *stateMachine) Lookup(q any) (any, error) {
	queryObj, ok := q.(query)
	if !ok {
		return nil, fmt.Errorf("invalid query type")
	}
	switch queryObj.Type {
	case "get":
		value, ok := s.kv[queryObj.Key]
		return getResult{Value: value, Found: ok}, nil
	case "ping":
		return true, nil
	default:
		return nil, nil
	}
}

func (s *stateMachine) SaveSnapshot(w io.Writer, _ sm.ISnapshotFileCollection, _ <-chan struct{}) error {
	enc := json.NewEncoder(w)
	return enc.Encode(s.kv)
}

func (s *stateMachine) RecoverFromSnapshot(r io.Reader, _ []sm.SnapshotFile, _ <-chan struct{}) error {
	kv := map[string]string{}
	dec := json.NewDecoder(r)
	if err := dec.Decode(&kv); err != nil {
		return err
	}
	s.kv = kv
	return nil
}

func (s *stateMachine) Close() error {
	return nil
}
