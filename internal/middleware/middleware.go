package middleware

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
)

// Cors provides CORS support.
func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET,POST")
		w.Header().Add("Access-Control-Allow-Headers", "content-type, authorization")

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GZIP(next http.Handler) http.Handler {
	return gziphandler.GzipHandler(next)
}
