package main

import (
	"net/http"

	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
	"github.com/stefanprodan/syros/models"
)

func (s *HttpServer) deploymentApiRoutes() chi.Router {
	r := chi.NewRouter()

	r.Post("/start", func(w http.ResponseWriter, r *http.Request) {
		d := Deployment{}
		if err := render.Bind(r, &d); err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			return
		}

		if len(d.TicketId) < 1 {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, "ticket_id is required")
		}

		if err := s.Repository.DeploymentStartUpsert(d.Deployment); err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			return
		}
	})

	r.Post("/finish", func(w http.ResponseWriter, r *http.Request) {
		d := Deployment{}
		if err := render.Bind(r, &d); err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			return
		}

		if len(d.TicketId) < 1 {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, "ticket_id is required")
		}

		if err := s.Repository.DeploymentUpsert(d.Deployment); err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			return
		}
	})
	return r
}

type Deployment struct {
	models.Deployment
}

func (l *Deployment) Bind(r *http.Request) error {
	// just a post-process after a decode..
	return nil
}
