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

	r.Post("/api/login", func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Username string `json:"name"`
			Password string `json:"password"`
		}

		if err := render.Bind(r.Body, &data); err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			return
		}

		if data.Username == "admin" && data.Password == "admin" {
			_, tokenString, err := s.TokenAuth.Encode(jwtauth.Claims{"username": data.Username})

			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}

			render.JSON(w, r, tokenString)

		} else {
			render.Status(r, http.StatusNotFound)
			render.PlainText(w, r, "Invalid Username or Password")
		}
	})

	r.Get("/api/error", func(w http.ResponseWriter, r *http.Request) {
		err := errors.New("This is just a test")
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, err.Error())
	})

	r.Group(func(r chi.Router) {
		r.Use(s.TokenAuth.Verifier)
		r.Use(jwtauth.Authenticator)

		r.Get("/api/hosts", func(w http.ResponseWriter, r *http.Request) {
			hosts, err := s.Repository.AllHosts()
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}

			render.JSON(w, r, hosts)
		})
	})

	http.ListenAndServe(fmt.Sprintf(":%v", s.Config.Port), r)
}
