package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

func containerHealthCheck(url string, timeout int) error {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(timeout) * time.Second,
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
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Print(string(body))

	if resp.StatusCode == 200 {
		return nil
	} else {
		return errors.Errorf("HTTP error code %d received", resp.StatusCode)
	}

	return nil
}
