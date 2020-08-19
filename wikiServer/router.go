package wikiServer

import "net/http"

type Router interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	HandleFunc(pattern string, handler func(ResponseWriter http.ResponseWriter, Request *http.Request))
	Handle(pattern string, handler http.Handler)
}
