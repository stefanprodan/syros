package main

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/braintree/manners"
	unrender "github.com/unrolled/render"
)

type HttpServer struct {
	Config   *Config
	Payloads map[string]*DockerPayload
}

// Starts HTTP Server
func (s *HttpServer) Start() {

	render := unrender.New(unrender.Options{
		IndentJSON: true,
		Layout:     "layout",
	})

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		render.JSON(w, http.StatusOK, s.Payloads)
	})
	http.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		render.Text(w, http.StatusOK, "pong")
	})
	http.HandleFunc("/config", func(w http.ResponseWriter, req *http.Request) {
		render.JSON(w, http.StatusOK, s.Config)
	})

	log.Fatal(manners.ListenAndServe(fmt.Sprintf(":%v", s.Config.Port), http.DefaultServeMux))
}

// Stop attempts to gracefully shutdown the HTTP server
func (s *HttpServer) Stop() {
	manners.Close()
}

