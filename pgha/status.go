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
	Code      int
	Message   string
	Timestamp time.Time
}

func NewStatus() *Status {
	return &Status{
		Code:      CandidateCode,
		Message:   "Initializing",
		Timestamp: time.Now().UTC(),
	}
}

func (s *Status) GetStatus() (int, string, time.Time) {
	var code int
	var msg string
	var ts time.Time
	s.Lock()
	defer s.Unlock()
	code = s.Code
	msg = s.Message
	ts = s.Timestamp
	return code, msg, ts
}

func (s *Status) SetStatus(code int, msg string) {
	s.Lock()
	defer s.Unlock()
	s.Code = code
	s.Message = msg
	s.Timestamp = time.Now().UTC()
}
