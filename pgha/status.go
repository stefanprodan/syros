package main

import (
	"sync"
	"time"
)

const (
	LeaderCode    = 200
	FollowerCode  = 409
	CandidateCode = 500
	FaultedCode   = 503
)

type Status struct {
	sync.Mutex
	Code             int
	Message          string
	Timestamp        time.Time
	IsConsulLeader   bool
	IsPostgresMaster bool
}

func NewStatus() *Status {
	return &Status{
		Code:      CandidateCode,
		Message:   "Initializing",
		Timestamp: time.Now().UTC(),
	}
}

func (s *Status) GetStatus() Status {
	s.Lock()
	defer s.Unlock()
	return *s
}

func (s *Status) SetConsulStatus(leader bool, code int, msg string) {
	s.Lock()
	defer s.Unlock()
	s.IsConsulLeader = leader
	s.Code = code
	s.Message = msg
	s.Timestamp = time.Now().UTC()
}

func (s *Status) SetPostgresStatus(leader bool) {
	s.Lock()
	defer s.Unlock()
	s.IsPostgresMaster = leader
}

func (s *Status) IsMaster() bool {
	s.Lock()
	defer s.Unlock()
	return s.IsPostgresMaster
}
