package main

import (
	"errors"
	_ "expvar"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	rnd "github.com/unrolled/render"
)

type HttpServer struct {
	config *Config
	status *Status
	election *Election
}

func NewHttpServer(config *Config, status *Status, election *Election) (*HttpServer, error) {
	if config.Port < 1 {
		return nil, errors.New("HTTP Server port is required")
	}

	server := &HttpServer{
		config: config,
		status: status,
		election: election,
	}

	return server, nil
}

// Starts HTTP Server
func (s *HttpServer) Start() {

	render := rnd.New(rnd.Options{
		IndentJSON: true,
		Layout:     "layout",
	})

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		render.Text(w, http.StatusOK, "pong")
	})
	http.HandleFunc("/config", func(w http.ResponseWriter, req *http.Request) {
		render.JSON(w, http.StatusOK, s.config)
	})
	http.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
		code, msg, ts := s.status.GetStatus()
		info := map[string]string{
			"status":    msg,
			"timestamp": ts.Format(time.RFC3339),
		}
		render.JSON(w, code, info)
	})
	http.HandleFunc("/version", func(w http.ResponseWriter, req *http.Request) {
		info := map[string]string{
			"syros_version": version,
			"os":            runtime.GOOS,
			"arch":          runtime.GOARCH,
			"golang":        runtime.Version(),
			"max_procs":     strconv.FormatInt(int64(runtime.GOMAXPROCS(0)), 10),
			"goroutines":    strconv.FormatInt(int64(runtime.NumGoroutine()), 10),
			"cpu_count":     strconv.FormatInt(int64(runtime.NumCPU()), 10),
		}
		render.JSON(w, http.StatusOK, info)
	})
	http.HandleFunc("/fallback", func(w http.ResponseWriter, req *http.Request) {
		err := s.election.Fallback()
		if err != nil{
			info := map[string]string{
				"status":    err.Error(),
			}
			render.JSON(w, http.StatusInternalServerError, info)
		}else {
			info := map[string]string{
				"status": "leadership lost",
			}
			render.JSON(w, http.StatusOK, info)
		}
	})

	log.Error(http.ListenAndServe(fmt.Sprintf(":%v", s.config.Port), http.DefaultServeMux))
}
