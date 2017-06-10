package models

import "time"

type ClusterPayload struct {
	HealthCheck ClusterHealthCheck `json:"health_check"`
	Environment string              `json:"environment"`
}

type ClusterHealthCheck struct {
	Id          string    `bson:"_id,omitempty" json:"id"`
	Status      string    `bson:"status" json:"status"`
	Output      string    `bson:"output" json:"output"`
	ServiceName string    `bson:"service_name" json:"service_name"`
	Collected   time.Time `bson:"collected" json:"collected"`
	Since       time.Time `bson:"since" json:"since"`
	Environment string    `bson:"environment" json:"environment"`
}

type ClusterHealthCheckLog struct {
	Id          string    `bson:"_id,omitempty" json:"id"`
	CheckId     string    `bson:"check_id,omitempty" json:"check_id"`
	Status      string    `bson:"status" json:"status"`
	Output      string    `bson:"output" json:"output"`
	ServiceName string    `bson:"service_name" json:"service_name"`
	Begin       time.Time `bson:"begin" json:"begin"`
	End         time.Time `bson:"end" json:"end"`
	Timestamp   time.Time `bson:"timestamp" json:"timestamp"`
	Duration    int64     `bson:"duration" json:"duration"`
	Environment string    `bson:"environment" json:"environment"`
}

func NewClusterHealthCheckLog(check ClusterHealthCheck, begin time.Time, end time.Time) ClusterHealthCheckLog {
	log := ClusterHealthCheckLog{
		Begin:       begin,
		End:         end,
		CheckId:     check.Id,
		Environment: check.Environment,
		Output:      check.Output,
		ServiceName: check.ServiceName,
		Status:      check.Status,
		Timestamp:   time.Now().UTC(),
	}

	log.Duration = int64(end.Sub(begin).Seconds())

	return log
}
