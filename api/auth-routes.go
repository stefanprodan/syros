package main

import (
	"github.com/goware/jwtauth"
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
	"net/http"
	"strings"
)

func (s *HttpServer) authRoutes() chi.Router {
	r := chi.NewRouter()

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		var data struct {
			Username string `json:"name"`
			Password string `json:"password"`
		}
		if err := render.Bind(r.Body, &data); err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			return
		}
		credentials := strings.Split(s.Config.Credentials, "@")
		if data.Username == credentials[0] && data.Password == credentials[1] {
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

	return r
}
