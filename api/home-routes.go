package main

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-chi/jwtauth"
	"github.com/stefanprodan/syros/models"
)

func (s *HttpServer) homeRoutes() chi.Router {
	r := chi.NewRouter()

	// JWT protected
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(s.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			dockerStats, err := s.Repository.EnvironmentHostContainerSum()
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}
			vsphere, err := s.Repository.AllVSphere()
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}

			hs := make([]models.VSphereHost, 0)
			for _, d := range vsphere.Hosts {
				d.Name = strings.Split(d.Name, ".")[0]
				if d.VMs > 0 {
					hs = append(hs, d)
				}
			}

			data := struct {
				DockerStats  []models.EnvironmentStats `json:"docker"`
				VSphereHosts []models.VSphereHost      `json:"vsphere"`
			}{
				DockerStats:  dockerStats,
				VSphereHosts: hs,
			}

			render.JSON(w, r, data)
		})

		r.Get("/syrosservices", func(w http.ResponseWriter, r *http.Request) {
			services, err := s.Repository.AllSyrosServices()
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}
			render.JSON(w, r, services)
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
