package auth

import (
	"log"
	"net/http"
)

var credentials = map[string]string{
	"secureworks": "supersecret",
}

func Basic() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if !ok {
				log.Println("No Credentails Provided")
				//TODO fix this
				//w.WriteHeader(http.StatusUnauthorized)
				next.ServeHTTP(w, r)
				return
			}
			if credentials[username] != password {
				log.Println("Invalid Credentails Provided")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
