package main

import (
	"github.com/goware/jwtauth"
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
	"github.com/stefanprodan/syros/models"
	"net/http"
)

func (s *HttpServer) releaseRoutes() chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(s.TokenAuth.Verifier)
		r.Use(jwtauth.Authenticator)

		r.Get("/all", func(w http.ResponseWriter, r *http.Request) {
			rels, err := s.Repository.AllReleases()
			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}

			chart := models.ChartDto{
				Labels: make([]string, 0),
				Values: make([]int64, 0),
			}
			deployments := 0
			// aggregate chart per day based on release end date
			for _, cont := range rels {
				deployments += cont.Deployments
				date := cont.End.Format("06-01-02")
				found := -1
				for i, s := range chart.Labels {
					if s == date {
						found = i
						break
					}
				}
				if found > -1 {
					chart.Values[found] += int64(cont.Deployments)
				} else {
					chart.Labels = append(chart.Labels, date)
					chart.Values = append(chart.Values, int64(cont.Deployments))
				}
			}

			data := struct {
				Releases    []models.Release `json:"releases"`
				Chart       models.ChartDto  `json:"chart"`
				Deployments int              `json:"deployments"`
			}{
				Releases:    rels,
				Chart:       chart,
				Deployments: deployments,
			}

			render.JSON(w, r, data)
		})
	})

	return r
}
