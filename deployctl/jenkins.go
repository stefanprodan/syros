package main

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type JenkinsConfig struct {
	API struct {
		Password string `yaml:"password"`
		URL      string `yaml:"url"`
		User     string `yaml:"user"`
	} `yaml:"api"`
}

func loadJenkinsConfig(dir string, name string) (JenkinsConfig, bool, error) {
	cfg := JenkinsConfig{}
	planPath := ""
	dir = path.Join(dir, "deploy", "integrations")
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, name+".yml") || strings.Contains(path, name+".yaml") {
			planPath = path
		}
		return nil
	})

	if err != nil {
		return cfg, false, errors.Wrapf(err, "Reading from %v failed", dir)
	}

	if len(planPath) < 1 {
		return cfg, false, nil
	}

	data, err := ioutil.ReadFile(planPath)
	if err != nil {
		return cfg, false, errors.Wrapf(err, "Reading %v failed", planPath)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, false, errors.Wrapf(err, "Parsing %v failed", planPath)
	}

	return cfg, true, nil
}
