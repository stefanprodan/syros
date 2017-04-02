package models

import "time"

type ConsulPayload struct {
	HealthChecks []ConsulHealthCheck `json:"health_checks"`
	Environment  string              `json:"environment"`
}

type ConsulHealthCheck struct {
	Id          string    `bson:"_id,omitempty" json:"id"`
	HostId      string    `bson:"host_id" json:"host_id"`
	HostName    string    `bson:"host_name" json:"host_name"`
	Node        string    `bson:"node" json:"node"`
	CheckID     string    `bson:"check_id" json:"check_id"`
	Name        string    `bson:"name" json:"name"`
	Status      string    `bson:"status" json:"status"`
	Notes       string    `bson:"notes" json:"notes"`
	Output      string    `bson:"output" json:"output"`
	ServiceID   string    `bson:"service_id" json:"service_id"`
	ServiceName string    `bson:"service_name" json:"service_name"`
	Collected   time.Time `bson:"collected" json:"collected"`
	Since       time.Time `bson:"since" json:"since"`
	Environment string    `bson:"environment" json:"environment"`
}
