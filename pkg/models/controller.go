package model

import (
	"encoding/json"
	"time"
)

type Agent struct {
	Id                  string    `gorm:"primaryKey;unique" json:"agent_id"`
	Name                string    `json:"name"`
	Host                string    `json:"host"`
	PollIntervalSeconds int       `json:"poll_interval_seconds"`
	CreatedAt           time.Time `json:"created_at"`
	LastSeen            time.Time `json:"last_seen"`
}

func (a *Agent) TableName() string {
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
	ID        uint            `gorm:"primaryKey;autoIncrement:true;column:id;unique" json:"-"`
	Version   int             `gorm:"index;column:version" json:"-"`
	Data      json.RawMessage `gorm:"column:data" json:"data" swaggertype:"object"`
	CreatedAt time.Time       `gorm:"column:created_at" json:"-"`
}

func (a *Configuration) TableName() string {
	return "configurations"
}
