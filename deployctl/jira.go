package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"bytes"
	"encoding/json"
	"log"

	"github.com/codeskyblue/go-sh"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type JiraConfig struct {
	API struct {
		Password string `yaml:"password"`
		URL      string `yaml:"url"`
		User     string `yaml:"user"`
	} `yaml:"api"`
}

func loadJiraConfig(dir string, name string) (JiraConfig, bool, error) {
	plan := JiraConfig{}
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

func (j JiraConfig) Post(ticket string, env string, component string, msg string) error {
	url := fmt.Sprintf("%s/issue/%s/comment", j.API.URL, ticket)
	log.Printf("Updateing Jira ticket %s", url)
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
		"body": fmt.Sprintf("Deployment of '%s' on environment '%s' host '%s' finished.", component, env, msg),
	}
	jsonData, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.Wrapf(err, "Jira request %s failed", url)
	}
	req.SetBasicAuth(j.API.User, j.API.Password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "Jira HTTP POST %s failed", url)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "Jira HTTP POST %s body read failed", url)
	}

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		return nil
	} else {
		log.Print(string(body))
		return errors.Errorf("Jira HTTP error code %d received", resp.StatusCode)
	}

	return nil

}

func (j JiraConfig) Upload(ticket string, dir string, file string) error {
	url := fmt.Sprintf("%s/issue/%s/attachments", j.API.URL, ticket)
	attachment := path.Join(dir, file)

	session := sh.NewSession()
	session.SetDir(dir)
	cmd := fmt.Sprintf("curl -D- -u %s:%s -X POST -H \"X-Atlassian-Token: nocheck\" -F \"file=@%s\" %s",
		j.API.User, j.API.Password, attachment, url)
	output, err := session.Command("/bin/sh", "-c", cmd).CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "Jira upload %s failed %s", url, output)
	}
	log.Printf("Deploy log uploaded to Jira ticket %s", url)
	return nil
}
