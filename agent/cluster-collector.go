package main

import (
	"io/ioutil"
	"net"
	"net/http"
	"runtime"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/stefanprodan/syros/models"
)

type ClusterCollector struct {
	Name        string
	ApiAddress  string
	Environment string
	Topic       string
}

func NewClusterCollector(name string, address string, env string) (*ClusterCollector, error) {
	c := &ClusterCollector{
		Name:        name,
		ApiAddress:  address,
		Environment: env,
		Topic:       "cluster",
	}

	return c, nil
}

func (col *ClusterCollector) Collect() (*models.ClusterPayload, error) {
	start := time.Now().UTC()
	client := &http.Client{
		Transport: DefaultTransport(),
	}

	payload := &models.ClusterPayload{
		Environment: col.Environment,
		HealthCheck: models.ClusterHealthCheck{
			ServiceName: col.Name,
			Environment: col.Environment,
			Id:          models.Hash(col.ApiAddress),
			Collected:   time.Now().UTC(),
			Status:      "offline",
		},
	}

	resp, err := client.Get(col.ApiAddress)
	if err != nil {
		payload.HealthCheck.Output = err.Error()
		return payload, err
	}

	if resp.StatusCode == 200 {
		payload.HealthCheck.Status = "leader"
	} else {
		payload.HealthCheck.Status = "follower"
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		payload.HealthCheck.Output = err.Error()
		return payload, nil
	}

	payload.HealthCheck.Output = string(body)

	log.Debugf("%v collect duration: %v status %v", col.ApiAddress, time.Now().UTC().Sub(start), payload.HealthCheck.Status)
	return payload, nil
}

func DefaultTransport() *http.Transport {
	transport := DefaultPooledTransport()
	transport.DisableKeepAlives = true
	transport.MaxIdleConnsPerHost = -1
	return transport
}

func DefaultPooledTransport() *http.Transport {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
	}
	return transport
}
