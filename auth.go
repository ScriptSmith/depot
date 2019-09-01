package main

import "net/http"

type Auth struct {
	user string
	pass string
}

func (auth *Auth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()

		if auth.user != "" || auth.pass != "" {
			if user == "" && pass == "" {
				w.Header().Set("WWW-Authenticate", "Basic realm=\"Depot username and password required\"")
				http.Error(w, "Unauthorized.", http.StatusUnauthorized)
				return
			} else if user != auth.user || pass != auth.pass {
				http.Error(w, "Unauthorized.", http.StatusUnauthorized)
				return
			} else {
				next.ServeHTTP(w, r)
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
