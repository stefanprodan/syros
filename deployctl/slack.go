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

type SlackConfig struct {
	API struct {
		Password string `yaml:"password"`
		URL      string `yaml:"url"`
		User     string `yaml:"user"`
	} `yaml:"api"`
}

func loadSlackConfig(dir string, name string) (SlackConfig, bool, error) {
	plan := SlackConfig{}
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

func (j SlackConfig) Post(ticket string, action string, env string, component string, host string) error {
	url := j.API.URL
	log.Printf("Updating Slack %s", url)
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
		"username":   "syros-deployctl",
		"icon_emoji": "rocket",
		"text":       fmt.Sprintf("%s of '%s' on environment '%s' host '%s' finished. JIRA ticket %s", action, component, env, host, ticket),
	}
	jsonData, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.Wrapf(err, "Slack request %s failed", url)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "Slack HTTP POST %s failed", url)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "Slack HTTP POST %s body read failed", url)
	}

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		return nil
	} else {
		log.Print(string(body))
		return errors.Errorf("Slack HTTP error code %d received", resp.StatusCode)
	}

	return nil

}
