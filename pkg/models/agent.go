package model

import (
	"encoding/json"
	"sync"
)

type AgentState struct {
	mu                  sync.RWMutex
	AgentID             string          `json:"agent_id"`
	ETag                string          `json:"etag"`
	PollUrl             string          `json:"poll_url"`
	PollIntervalSeconds int             `json:"poll_interval_seconds"`
	Config              json.RawMessage `json:"config"`
}

func (s *AgentState) RegistraionData(agentID, pollUrl string, interval int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.AgentID = agentID
	s.PollUrl = pollUrl
	s.PollIntervalSeconds = interval
}

func (s *AgentState) UpdateConfig(etag string, config []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ETag = etag
	s.Config = config
}

func (s *AgentState) Get() (string, string, string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.AgentID, s.ETag, s.PollUrl
}

func (s *AgentState) GetInterval() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.PollIntervalSeconds <= 0 {
		return 30
	}

	return s.PollIntervalSeconds
}

func (s *AgentState) Snapshot() *AgentState {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return &AgentState{
		AgentID:             s.AgentID,
		ETag:                s.ETag,
		PollUrl:             s.PollUrl,
		PollIntervalSeconds: s.PollIntervalSeconds,
		Config:              s.Config,
	}
}

type ConfigResponse struct {
	ETag string
	Data json.RawMessage `json:"data"`
}
