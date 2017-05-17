package main

import (
	"net/http"

	"github.com/goware/jwtauth"
	"github.com/stefanprodan/chi"
	"github.com/stefanprodan/chi/render"
)

func (s *HttpServer) vsphereRoutes() chi.Router {
	r := chi.NewRouter()

	// JWT protected
	r.Group(func(r chi.Router) {
		r.Use(s.TokenAuth.Verifier)
		r.Use(jwtauth.Authenticator)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			vsphere, err := s.Repository.AllVSphere()
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}
			render.JSON(w, r, vsphere)
		})

	})

	return r
}
