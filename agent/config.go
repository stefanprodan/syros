package main

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Config holds global configuration, defaults are provided in main.
// Agent config is populated from startup flag.
type Config struct {
	Environment     string `m:"Environment"`
	LogLevel        string `m:"LogLevel"`
	Port            int    `m:"Port"`
	Nats            string `m:"Nats"`
	CollectorConfig string `m:"CollectorConfig"`
}

type CollectorConfig struct {
	Docker  ApiCollectorConfig `json:"docker" yaml:"docker"`
	Consul  ApiCollectorConfig `json:"consul" yaml:"consul"`
	VSphere ApiCollectorConfig `json:"vsphere" yaml:"vsphere"`
}

type ApiCollectorConfig struct {
	Endpoints []string `json:"endpoints" yaml:"endpoints"`
	Include   []string `json:"include" yaml:"include"`
	Exclude   []string `json:"exclude" yaml:"exclude"`
	Cron      string   `json:"cron" yaml:"cron"`
}

func LoadCollectorConfig(path string) (*CollectorConfig, error) {
	cfg := &CollectorConfig{}

	if len(path) < 1 {
		return nil, errors.Errorf("Collector config %v not found", path)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "Reading %v failed", path)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, errors.Wrapf(err, "Parsing %v failed", path)
	}

	return cfg, nil
}
