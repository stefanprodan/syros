package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Config holds global configuration, defaults are provided in main.
// Agent config is populated from startup flag.
type Config struct {
	Environment            string `m:"Environment"`
	LogLevel               string `m:"LogLevel"`
	Port                   int    `m:"Port"`
	CollectInterval        int    `m:"CollectInterval"`
	DockerApiAddresses     string `m:"DockerApiAddresses"`
	ConsulApiAddresses     string `m:"ConsulApiAddresses"`
	VSphereApiAddress      string `m:"VSphereApiAddresses"`
	VSphereInclude         string `m:"VSphereInclude"`
	VSphereExclude         string `m:"VSphereExclude"`
	VSphereCollectInterval int    `m:"VSphereCollectInterval"`
	Nats                   string `m:"Nats"`
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

func LoadCollectorConfig(dir string, name string) (*CollectorConfig, error) {
	cfg := &CollectorConfig{}
	cfgPath := ""
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, name+".yml") || strings.Contains(path, name+".yaml") {
			cfgPath = path
		}
		return nil
	})

	if err != nil {
		return nil, errors.Wrapf(err, "Reading from %v failed", dir)
	}

	if len(cfgPath) < 1 {
		return nil, errors.Errorf("Collector config %v not found", name)
	}

	data, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Reading %v failed", cfgPath)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, errors.Wrapf(err, "Parsing %v failed", cfgPath)
	}

	return cfg, nil
}
