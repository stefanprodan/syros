package main

import (
	"fmt"
	"github.com/goware/cors"
	"github.com/goware/jwtauth"
	"github.com/pkg/errors"
	"github.com/stefanprodan/chi"
	"github.com/stefanprodan/chi/middleware"
	"github.com/stefanprodan/chi/render"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
)

type HttpServer struct {
	Config     *Config
	Repository *Repository
	TokenAuth  *jwtauth.JwtAuth
}

func (s *HttpServer) Start() {

	prom := NewPrometheus("syros", "api", true, true)

	r := chi.NewRouter()

	corsWare := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(corsWare.Handler)
	r.Use(prom.Middleware)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.DefaultCompress)
	if s.Config.LogLevel == "debug" {
		r.Use(middleware.DefaultLogger)
	}

	r.Mount("/metrics", prom.Router())
	r.Mount("/debug", s.pprofRoutes())

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, "pong")
	})

	r.Get("/config", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, s.Config)
	})

	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, "OK")
	})

	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		info := map[string]string{
			"syros_version": version,
			"os":            runtime.GOOS,
			"arch":          runtime.GOARCH,
			"golang":        runtime.Version(),
			"max_procs":     strconv.FormatInt(int64(runtime.GOMAXPROCS(0)), 10),
			"goroutines":    strconv.FormatInt(int64(runtime.NumGoroutine()), 10),
			"cpu_count":     strconv.FormatInt(int64(runtime.NumCPU()), 10),
		}

		render.JSON(w, r, info)
	})

	r.Get("/api/error", func(w http.ResponseWriter, r *http.Request) {
		err := errors.New("This is just a test")
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
	})

	r.Get("/api/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("This is just a test")
	})

	r.Mount("/api/auth", s.authRoutes())
	r.Mount("/api/docker", s.dockerRoutes())
	r.Mount("/api/deployment", s.deploymentApiRoutes())
	r.Mount("/api/release", s.releaseRoutes())

	// ui paths
	indexPath := filepath.Join(s.Config.AppPath, "index.html")
	staticPath := filepath.Join(s.Config.AppPath, "static")

	// set index.html as entry point
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, indexPath)
	})

	// static files (js, css, fonts)
	r.FileServer("/static", http.Dir(staticPath))

	err := http.ListenAndServe(fmt.Sprintf(":%v", s.Config.Port), r)
	if err != nil {
		log.Fatalf("HTTP Server crashed! %v", err.Error())
	}
}
