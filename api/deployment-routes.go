package main

import (
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
	"github.com/stefanprodan/syros/models"
	"net/http"
)

func (s *HttpServer) deploymentApiRoutes() chi.Router {
	r := chi.NewRouter()

	r.Post("/begin", func(w http.ResponseWriter, r *http.Request) {
		d := models.Deployment{}
		if err := render.Bind(r.Body, &d); err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			return
		}
	})

	r.Post("/end", func(w http.ResponseWriter, r *http.Request) {
		d := models.Deployment{}
		if err := render.Bind(r.Body, &d); err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			return
		}
	})
	return r
}
