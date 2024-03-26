package middlewares

import (
	"encoding/json"
	"net/http"
	"os"
)

type Config struct {
	Users []User `json:"users"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func BasicAuthMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
			config, err := loadConfig("users.json")
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			user, pass, ok := rq.BasicAuth()
			if !ok {
				unauthorized(w)
				return
			}

			if !authenticate(user, pass, config.Users) {
				unauthorized(w)
				return
			}

			next.ServeHTTP(w, rq)
		})
	}
}

func loadConfig(configFile string) (Config, error) {
	var config Config
	file, err := os.Open(configFile)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}

func authenticate(username, password string, users []User) bool {
	for _, user := range users {
		if username == user.Username && password == user.Password {
			return true
		}
	}
	return false
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}
