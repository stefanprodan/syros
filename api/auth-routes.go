package main

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-chi/jwtauth"
)

func (s *HttpServer) authRoutes() chi.Router {
	r := chi.NewRouter()

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		data := LoginForm{}
		if err := render.Bind(r, &data); err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.PlainText(w, r, err.Error())
			return
		}
		credentials := strings.Split(s.Config.Credentials, "@")
		if data.Username == credentials[0] && data.Password == credentials[1] {
			var claims = jwtauth.Claims{"role": data.Username}
			claims = claims.SetIssuedNow()
			// TODO: set expiry based on role
			// claims = claims.SetExpiry(time.Now().Add(time.Hour * 48))
			_, tokenString, err := s.TokenAuth.Encode(claims)
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

type LoginForm struct {
	Username string `json:"name"`
	Password string `json:"password"`
}

func (l *LoginForm) Bind(r *http.Request) error {
	// just a post-process after a decode..
	return nil
}
