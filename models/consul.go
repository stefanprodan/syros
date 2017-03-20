package models

import "time"

type ConsulPayload struct {
	HealthChecks []ConsulHealthCheck `json:"health_checks"`
	Environment  string              `json:"environment"`
}

type ConsulHealthCheck struct {
	Id          string    `gorethink:"id,omitempty" json:"id"`
	HostId      string    `gorethink:"host_id" json:"host_id"`
	HostName    string    `gorethink:"host_name" json:"host_name"`
	Node        string    `gorethink:"node" json:"node"`
	CheckID     string    `gorethink:"check_id" json:"check_id"`
	Name        string    `gorethink:"name" json:"name"`
	Status      string    `gorethink:"status" json:"status"`
	Notes       string    `gorethink:"notes" json:"notes"`
	Output      string    `gorethink:"output" json:"output"`
	ServiceID   string    `gorethink:"service_id" json:"service_id"`
	ServiceName string    `gorethink:"service_name" json:"service_name"`
	Collected   time.Time `gorethink:"collected" json:"collected"`
	Since       time.Time `gorethink:"since" json:"since"`
	Environment string    `gorethink:"environment" json:"environment"`
}
