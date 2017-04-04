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

type ConsulHealthCheckLog struct {
	Id          string    `bson:"_id,omitempty" json:"id"`
	CheckId     string    `bson:"check_id,omitempty" json:"check_id"`
	HostId      string    `bson:"host_id" json:"host_id"`
	HostName    string    `bson:"host_name" json:"host_name"`
	Node        string    `bson:"node" json:"node"`
	Name        string    `bson:"name" json:"name"`
	Status      string    `bson:"status" json:"status"`
	Notes       string    `bson:"notes" json:"notes"`
	Output      string    `bson:"output" json:"output"`
	ServiceID   string    `bson:"service_id" json:"service_id"`
	ServiceName string    `bson:"service_name" json:"service_name"`
	Begin       time.Time `bson:"begin" json:"begin"`
	End         time.Time `bson:"end" json:"end"`
	Timestamp   time.Time `bson:"timestamp" json:"timestamp"`
	Environment string    `bson:"environment" json:"environment"`
}

func NewConsulHealthCheckLog(check ConsulHealthCheck, begin time.Time, end time.Time) ConsulHealthCheckLog {
	log := ConsulHealthCheckLog{
		Begin:       begin,
		End:         end,
		CheckId:     check.Id,
		Environment: check.Environment,
		HostId:      check.HostId,
		HostName:    check.HostName,
		Name:        check.Name,
		Node:        check.Node,
		Notes:       check.Notes,
		Output:      check.Output,
		ServiceID:   check.ServiceID,
		ServiceName: check.ServiceName,
		Status:      check.Status,
		Timestamp:   time.Now().UTC(),
	}

	return log
}
