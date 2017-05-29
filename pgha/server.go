package main

import (
	_ "expvar"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	rnd "github.com/unrolled/render"
)

type HttpServer struct {
	Config *Config
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
		render.JSON(w, http.StatusOK, s.Config)
	})
	http.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
		render.Text(w, http.StatusOK, "OK")
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

	log.Error(http.ListenAndServe(fmt.Sprintf(":%v", s.Config.Port), http.DefaultServeMux))
}
