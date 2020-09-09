package middlewares

import (
	"fmt"
	"net/http"
)

func JsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")

		fmt.Println("entrou")

		next.ServeHTTP(w, r)
		fmt.Println("saiu")
	})
}
