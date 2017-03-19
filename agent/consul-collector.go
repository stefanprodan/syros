package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	consul "github.com/hashicorp/consul/api"
	"github.com/stefanprodan/syros/models"
	"time"
)

type ConsulCollector struct {
	ApiAddress  string
	Environment string
	Client      *consul.Client
}

func NewConsulCollector(address string, env string) (*ConsulCollector, error) {
	cfg := consul.DefaultConfig()
	cfg.Address = address

	client, err := consul.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	c := &ConsulCollector{
		ApiAddress:  address,
		Environment: env,
		Client:      client,
	}

	return c, nil
}

func (col *ConsulCollector) Collect() (*models.ConsulPayload, error) {
	start := time.Now().UTC()

	health := col.Client.Health()

	checks, meta, err := health.State("any", nil)
	if err != nil {
		return nil, err
	}
	if meta.LastIndex == 0 {
		return nil, fmt.Errorf("Consul health check bad meta: %v", meta)
	}
	if len(checks) == 0 {
		return nil, fmt.Errorf("Consul no health checks found: %v", checks)
	}
	res := make([]models.ConsulHealthCheck, 0)
	for _, ck := range checks {
		res = append(res, MapConsulCheck(ck))
	}

	payload := &models.ConsulPayload{
		HealthChecks: res,
		Environment:  col.Environment,
	}

	log.Debugf("%v collect duration: %v health checks %v", col.ApiAddress, time.Now().UTC().Sub(start), len(payload.HealthChecks))
	return payload, nil
}

func MapConsulCheck(ck *consul.HealthCheck) models.ConsulHealthCheck {
	check := models.ConsulHealthCheck{
		CheckID:     ck.CheckID,
		Name:        ck.Name,
		Node:        ck.Node,
		Notes:       ck.Notes,
		Output:      ck.Output,
		ServiceID:   ck.ServiceID,
		ServiceName: ck.ServiceName,
		Status:      ck.Status,
	}
	return check
}
