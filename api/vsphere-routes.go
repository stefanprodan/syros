package main

import (
	"net/http"

	"github.com/goware/jwtauth"
	"github.com/stefanprodan/chi"
	"github.com/stefanprodan/chi/render"
	"github.com/stefanprodan/syros/models"
)

func (s *HttpServer) vsphereRoutes() chi.Router {
	r := chi.NewRouter()

	// JWT protected
	r.Group(func(r chi.Router) {
		r.Use(s.TokenAuth.Verifier)
		r.Use(jwtauth.Authenticator)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			vsphere, err := s.Repository.AllVSphere()

			chart := models.ChartDto{
				Labels: make([]string, 0),
				Values: make([]int64, 0),
			}

			for _, vm := range vsphere.VMs {
				found := -1
				for i, s := range chart.Labels {
					if s == vm.Environment {
						found = i
						break
					}
				}
				if found > -1 {
					chart.Values[found]++
				} else {
					chart.Labels = append(chart.Labels, vm.Environment)
					chart.Values = append(chart.Values, int64(1))
				}
			}

			data := struct {
				Hosts      []models.VSphereHost      `json:"hosts"`
				DataStores []models.VSphereDatastore `json:"data_stores"`
				VMs        []models.VSphereVM        `json:"vms"`
				Chart       models.ChartDto  `json:"chart"`
			}{
				Hosts:      vsphere.Hosts,
				DataStores: vsphere.DataStores,
				VMs:        vsphere.VMs,
				Chart:      chart,
			}

			if err != nil {
				render.Status(r, http.StatusInternalServerError)
				render.PlainText(w, r, err.Error())
				return
			}
			render.JSON(w, r, data)
		})

	})

	return r
}
