package auth

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var basicAuthCredentials credentials

func init() {
	credentialsFile, err := os.Open("credentials.json")
	if err != nil {
		log.Fatal("Error opening credentials file:", err.Error())
	}
	defer credentialsFile.Close()
	credData, err := ioutil.ReadAll(credentialsFile)
	if err != nil {
		log.Fatal("Error reading credentials file:", err.Error())
	}
	basicAuthCredentials = credentials{}
	if err := json.Unmarshal(credData, &basicAuthCredentials); err != nil {
		log.Fatal("Error parsing json in credentials file:", err.Error())
	}
}

func Basic() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if !ok {
				log.Println("No Credentails Provided")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if basicAuthCredentials.Username != username {
				log.Println("Invalid Username")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if basicAuthCredentials.Password != password {
				log.Println("Invalid Password")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
