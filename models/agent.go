package models

import "time"

type Agent struct {
	Id          string `gorethink:"id,omitempty"`
	Hostname    string
	LastSeen    time.Time
	Environment string
}
