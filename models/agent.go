package models

import "time"

type Agent struct {
	Id          string    `gorethink:"id,omitempty" json:"id"`
	Hostname    string    `gorethink:"hostname" json:"hostname"`
	LastSeen    time.Time `gorethink:"hostname" json:"last_seen"`
	Environment string    `gorethink:"environment" json:"environment"`
}
