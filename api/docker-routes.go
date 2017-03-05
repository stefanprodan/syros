package main

import (
	"github.com/goware/jwtauth"
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
	"net/http"
)

func (s *HttpServer) dockerRoutes() chi.Router {
	r := chi.NewRouter()

	// JWT protected
	r.Group(func(r chi.Router) {
		r.Use(s.TokenAuth.Verifier)
		r.Use(jwtauth.Authenticator)

		r.Get("/hosts", func(w http.ResponseWriter, r *http.Request) {
			hosts, err := s.Repository.AllHosts()
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}
			render.JSON(w, r, hosts)
		})

		r.Get("/hosts/:hostID", func(w http.ResponseWriter, r *http.Request) {
			hostID := chi.URLParam(r, "hostID")

			payload, err := s.Repository.HostContainers(hostID)
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}
			render.JSON(w, r, payload)
		})

		r.Get("/environment/:env", func(w http.ResponseWriter, r *http.Request) {
			env := chi.URLParam(r, "env")

			payload, err := s.Repository.EnvironmentContainers(env)
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}
			render.JSON(w, r, payload)
		})

		r.Get("/containers", func(w http.ResponseWriter, r *http.Request) {
			containers, err := s.Repository.AllContainers()
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}
			render.JSON(w, r, containers)
		})

		r.Get("/containers/:containerID", func(w http.ResponseWriter, r *http.Request) {
			containerID := chi.URLParam(r, "containerID")

			payload, err := s.Repository.Container(containerID)
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}
			render.JSON(w, r, payload)
		})
	})

	r.Get("/environments", func(w http.ResponseWriter, r *http.Request) {
		environments, err := s.Repository.AllEnvironments()
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			return
		}
		render.JSON(w, r, environments)
	})

	return r
}
