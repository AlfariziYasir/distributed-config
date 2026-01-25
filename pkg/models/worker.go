package model

import (
	"encoding/json"
	"sync"
)

type DataConfig struct {
	mu     sync.RWMutex
	Config json.RawMessage `json:"config"`
}

func (d *DataConfig) UpdateData(config json.RawMessage) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Config = config
}

func (d *DataConfig) GetConfig() json.RawMessage {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.Config
}
