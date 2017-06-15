package main

import (
	"net/http"

	"github.com/goware/jwtauth"
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
	"github.com/stefanprodan/syros/models"
)

func (s *HttpServer) clusterRoutes() chi.Router {
	r := chi.NewRouter()

	// JWT protected
	r.Group(func(r chi.Router) {
		r.Use(s.TokenAuth.Verifier)
		r.Use(jwtauth.Authenticator)

		r.Get("/healthchecks", func(w http.ResponseWriter, r *http.Request) {
			checks, err := s.Repository.AllClusterHealthChecks()
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}
			render.JSON(w, r, checks)
		})

		r.Get("/healthchecks/:checkID", func(w http.ResponseWriter, r *http.Request) {
			checkID := chi.URLParam(r, "checkID")
			checks, stats, err := s.Repository.ClusterHealthCheckLog(checkID)
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}

			data := struct {
				Checks []models.ClusterHealthCheckLog `json:"checks"`
				Stats  []models.HealthCheckStats      `json:"stats"`
			}{
				Checks: checks,
				Stats:  stats,
			}

			render.JSON(w, r, data)
		})

	})

	return r
}
