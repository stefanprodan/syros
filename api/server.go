package main

import (
	"fmt"
	"github.com/goware/cors"
	"github.com/goware/jwtauth"
	"github.com/pkg/errors"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/pressly/chi/render"
	"net/http"
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
	if s.Config.LogLevel == "debug" {
		r.Use(middleware.DefaultLogger)
	}

	r.Get("/api/error", func(w http.ResponseWriter, r *http.Request) {
		err := errors.New("This is just a test")
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
	})

	r.Mount("/api/auth", s.authRoutes())
	r.Mount("/api/docker", s.dockerRoutes())

	http.ListenAndServe(fmt.Sprintf(":%v", s.Config.Port), r)
}
