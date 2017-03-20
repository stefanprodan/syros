package main

import (
	"github.com/pkg/errors"
	"sync"
	"time"
)

type CollectorStatus struct {
	LastCollectTime    time.Time
	LastCollectSuccess bool
}

type AgentStatus struct {
	Collectors map[string]CollectorStatus
	mutex      sync.RWMutex
}

func NewAgentStatus(endpoints []string) (*AgentStatus, error) {
	if len(endpoints) < 1 {
		return nil, errors.New("No endpoints specified")
	}
	collectors := make(map[string]CollectorStatus, len(endpoints))
	for _, endpoint := range endpoints {
		collectors[endpoint] = CollectorStatus{
			LastCollectSuccess: false,
			LastCollectTime:    time.Now(),
		}
	}
	status := &AgentStatus{
		Collectors: collectors,
	}

	return status, nil
}

func (a *AgentStatus) SetCollectorStatus(endpoint string, lastCollectSuccess bool) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.Collectors[endpoint] = CollectorStatus{
		LastCollectSuccess: lastCollectSuccess,
		LastCollectTime:    time.Now(),
	}
}

func (a *AgentStatus) GetStatus() (*AgentStatus, int) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	status := &AgentStatus{
		Collectors: make(map[string]CollectorStatus, len(a.Collectors)),
	}
	code := 200
	for i, collector := range a.Collectors {
		status.Collectors[i] = CollectorStatus{
			LastCollectSuccess: collector.LastCollectSuccess,
			LastCollectTime:    collector.LastCollectTime,
		}
		if !collector.LastCollectSuccess {
			code = 500
		}
	}

	return status, code
}
