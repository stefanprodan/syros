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

type ComponentConfig struct {
	Component struct {
		Image  string            `yaml:"image"`
		Mode   string            `yaml:"mode"`
		Name   string            `yaml:"name"`
		Target []ComponentTarget `yaml:"target"`
		Type   string            `yaml:"type"`
	} `yaml:"component"`
}

type ComponentTarget struct {
	Dir        string `yaml:"dir"`
	Health     string `yaml:"health"`
	Host       string `yaml:"host"`
	Leadership string `yaml:"leadership"`
	DBName     string `yaml:"dbname"`
	DBUser     string `yaml:"user"`
	DBPassword string `yaml:"password"`
	URL        string `yaml:"url"`
}

func removeTargetByIndex(s []ComponentTarget, index int) []ComponentTarget {
	return append(s[:index], s[index+1:]...)
}

func loadComponent(dir string, env string, name string) (ComponentConfig, bool, error) {
	plan := ComponentConfig{}
	planPath := ""
	dir = path.Join(dir, "deploy", "components", env)
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, name+".yml") || strings.Contains(path, name+".yaml") {
			planPath = path
		}
		return nil
	})

	if err != nil {
		return plan, false, errors.Wrapf(err, "Reading from %v failed", dir)
	}

	if len(planPath) < 1 {
		return plan, false, nil
	}

	data, err := ioutil.ReadFile(planPath)
	if err != nil {
		return plan, false, errors.Wrapf(err, "Reading %v failed", planPath)
	}

	if err := yaml.Unmarshal(data, &plan); err != nil {
		return plan, false, errors.Wrapf(err, "Parsing %v failed", planPath)
	}

	return plan, true, nil
}

type PromotionConfig struct {
	Rules struct {
		Action string `yaml:"action"`
		Source string `yaml:"source"`
	} `yaml:"rules"`
}

func loadPromotion(dir string, name string) (PromotionConfig, error) {
	plan := PromotionConfig{}
	planPath := ""
	dir = path.Join(dir, "deploy", "defaults")
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, name+".yml") || strings.Contains(path, name+".yaml") {
			planPath = path
		}
		return nil
	})

	if err != nil {
		return plan, errors.Wrapf(err, "Reading from %v failed", dir)
	}

	if len(planPath) < 1 {
		return plan, errors.Errorf("Config %v not found", name)
	}

	data, err := ioutil.ReadFile(planPath)
	if err != nil {
		return plan, errors.Wrapf(err, "Reading %v failed", planPath)
	}

	if err := yaml.Unmarshal(data, &plan); err != nil {
		return plan, errors.Wrapf(err, "Parsing %v failed", planPath)
	}

	return plan, nil
}
