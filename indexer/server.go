package main

import (
	"fmt"
	"github.com/braintree/manners"
	unrender "github.com/unrolled/render"
	"net/http"
)

type HttpServer struct {
	Config     *Config
	Registry   *Registry
	Repository *Repository
}

// Starts HTTP Server
func (s *HttpServer) Start() {

	render := unrender.New(unrender.Options{
		IndentJSON: true,
		Layout:     "layout",
	})

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		render.JSON(w, http.StatusOK, s.Registry.GetActiveAgents())
	})
	http.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		render.Text(w, http.StatusOK, "pong")
	})
	http.HandleFunc("/config", func(w http.ResponseWriter, req *http.Request) {
		render.JSON(w, http.StatusOK, s.Config)
	})
	http.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
		render.JSON(w, http.StatusOK, s.Registry.GetActiveAgents())
	})

	manners.ListenAndServe(fmt.Sprintf(":%v", s.Config.Port), http.DefaultServeMux)
}

// Stop attempts to gracefully shutdown the HTTP server
func (s *HttpServer) Stop() {
	manners.Close()
}
