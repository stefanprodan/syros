package main

import (
	"fmt"
	"github.com/goware/cors"
	"github.com/goware/jwtauth"
	"github.com/pkg/errors"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/pressly/chi/render"
	"log"
	"net/http"
	"path/filepath"
)

type HttpServer struct {
	Config     *Config
	Repository *Repository
	TokenAuth  *jwtauth.JwtAuth
}

func (s *HttpServer) Start() {

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
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.DefaultCompress)
	if s.Config.LogLevel == "debug" {
		r.Use(middleware.DefaultLogger)
	}

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, "pong")
	})

	r.Get("/config", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, s.Config)
	})

	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, s.Config)
	})

	r.Get("/api/error", func(w http.ResponseWriter, r *http.Request) {
		err := errors.New("This is just a test")
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
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
