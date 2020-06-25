package main

import "net/http"

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func basic(username, password string, next http.Handler) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		u, pw, ok := req.BasicAuth()
		if !ok || u != username || pw != password {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, req)
	}
}
