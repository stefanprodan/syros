package main

import (
	"fmt"
	"github.com/braintree/manners"
	"github.com/pkg/errors"
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
	http.HandleFunc("/api/hosts", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		hosts, err := s.Repository.AllHosts()
		if err != nil {
			render.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		render.JSON(w, http.StatusOK, hosts)
	})
	http.HandleFunc("/api/containers", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		containers, err := s.Repository.AllContainers()
		if err != nil {
			render.JSON(w, http.StatusInternalServerError, err.Error())
			return
		}

		render.JSON(w, http.StatusOK, containers)
	})

	http.HandleFunc("/api/error", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		err := errors.New("This is just a test")
		render.Text(w, http.StatusInternalServerError, err.Error())
	})

	manners.ListenAndServe(fmt.Sprintf(":%v", s.Config.Port), http.DefaultServeMux)
}

// Stop attempts to gracefully shutdown the HTTP server
func (s *HttpServer) Stop() {
	manners.Close()
}
