package middlewares

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

func Logging() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				r.Body = io.NopCloser(bytes.NewBuffer(body))
				log.Printf("[%s] %s %s %s\n%s\n", time.Now().Format(time.DateTime), r.Method, r.URL.Path, r.RemoteAddr, string(body))
			}
			next.ServeHTTP(w, r)
		})
	}
}
