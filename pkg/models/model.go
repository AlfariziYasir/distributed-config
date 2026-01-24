package model

import (
	"encoding/json"
	"time"
)

type Agent struct {
	Id        string    `gorm:"primaryKey" json:"agent_id"`
	Name      string    `json:"name"`
	Host      string    `json:"host"`
	CreatedAt time.Time `json:"created_at"`
	LastSeen  time.Time `json:"last_seen"`
}

func (a *Agent) Tablename() string {
	return "agents"
}

type AgentRequest struct {
	Name string `json:"name"`
	Host string `json:"host"`
}

type AgentResponse struct {
	AgentId             string `json:"agent_id"`
	PollUrl             string `json:"poll_url"`
	PollIntervalSeconds int    `json:"poll_interval_seconds"`
}

type Configuration struct {
	ID        uint            `gorm:"primaryKey;autoIncrement:" json:"-"`
	Version   int             `gorm:"index" json:"version"`
	Data      json.RawMessage `json:"data"`
	CreatedAt time.Time       `json:"created_at"`
}

func (a *Configuration) Tablename() string {
	return "configurations"
}
