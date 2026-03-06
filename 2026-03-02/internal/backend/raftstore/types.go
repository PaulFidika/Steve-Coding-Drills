package raftstore

import "encoding/json"

type command struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type setCommand struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type query struct {
	Type string `json:"type"`
	Key  string `json:"key,omitempty"`
}

type getResult struct {
	Value string `json:"value"`
	Found bool   `json:"found"`
}
