package main

import (
	"fmt"
	"net/http"
	"regexp"
)

type router struct {
}

func (rou router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	for reg, handler := range ROUTER {
		if reg.MatchString(p) {
			w.Header().Add("Access-Control-Allow-Origin", "*")
			handler(w, r)
			return
		}
	}

	http.NotFound(w, r)
}

func NewRouter() router {
	for r, h := range myHandlers {
		ROUTER[regexp.MustCompile(fmt.Sprintf("%s$", r))] = h
	}
	return router{}
}
