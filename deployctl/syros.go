package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type SyrosConfig struct {
	API struct {
		Password string `yaml:"password"`
		URL      string `yaml:"url"`
		User     string `yaml:"user"`
	} `yaml:"api"`
}

func loadSyrosConfig(dir string, name string) (SyrosConfig, bool, error) {
	plan := SyrosConfig{}
	planPath := ""
	dir = path.Join(dir, "deploy", "integrations")
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

func (j SyrosConfig) Start(ticket string, env string, component string, host string) error {
	url := fmt.Sprintf("%s/deployment/start", j.API.URL)
	log.Printf("Updateing Syros %s", url)
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   15 * time.Second,
			KeepAlive: 10 * time.Second,
		}).DialContext,
		DisableKeepAlives:     true,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   -1,
	}

	client := &http.Client{
		Transport: transport,
	}

	data := map[string]string{
		"ticket_id":    ticket,
		"service_name": component,
		"host_name":    strings.Split(host, ".")[0],
		"environment":  env,
	}
	jsonData, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.Wrapf(err, "Syros request %s failed", url)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "Syros HTTP POST %s failed", url)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "Syros HTTP POST %s body read failed", url)
	}

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		return nil
	} else {
		log.Print(string(body))
		return errors.Errorf("Syros HTTP error code %d received", resp.StatusCode)
	}

	return nil

}

func (j SyrosConfig) Finish(ticket string, env string, component string, host string, logFile string) error {

	logData, err := ioutil.ReadFile(logFile)
	if err != nil {
		return errors.Wrapf(err, "Syros log read %s failed", logFile)
	}

	url := fmt.Sprintf("%s/deployment/finish", j.API.URL)
	log.Printf("Updateing Syros %s", url)
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   15 * time.Second,
			KeepAlive: 10 * time.Second,
		}).DialContext,
		DisableKeepAlives:     true,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   -1,
	}

	client := &http.Client{
		Transport: transport,
	}

	data := map[string]string{
		"ticket_id":    ticket,
		"service_name": component,
		"host_name":    strings.Split(host, ".")[0],
		"environment":  env,
		"log":          string(logData),
	}
	jsonData, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.Wrapf(err, "Syros request %s failed", url)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "Syros HTTP POST %s failed", url)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "Syros HTTP POST %s body read failed", url)
	}

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		return nil
	} else {
		log.Print(string(body))
		return errors.Errorf("Syros HTTP error code %d received", resp.StatusCode)
	}

	return nil

}
