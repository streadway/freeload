package freeload

import (
	"log"
	"net/http"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.Header.Get("Origin"), r.URL)
		next.ServeHTTP(w, r)
	})
}
