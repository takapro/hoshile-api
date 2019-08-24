package main

import (
	"net/http"
	"strings"
)

type Router struct {
	entries []entry
}

type Handlers = map[string]func(http.ResponseWriter, *http.Request)

type entry struct {
	path     string
	exact    bool
	handlers Handlers
}

func NewRouter() Router {
	return Router{[]entry{}}
}

func (r *Router) Handle(path string, handlers Handlers) {
	r.entries = append(r.entries, entry{path, !strings.HasSuffix(path, "/"), handlers})
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, e := range router.entries {
		if (e.exact && r.URL.Path == e.path) || (!e.exact && strings.HasPrefix(r.URL.Path, e.path)) {
			handler := e.handlers[r.Method]
			if handler != nil {
				handler(w, r)
				return
			}

			if r.Method == http.MethodOptions {
				handleOptions(w, e)
				return
			}

			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
	}

	http.Error(w, "Not found", http.StatusNotFound)
}

func handleOptions(w http.ResponseWriter, e entry) {
	methods := ""
	for method := range e.handlers {
		methods += method + ", "
	}
	methods += http.MethodOptions

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", methods)
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	w.Header().Set("Access-Control-Max-Age", "86400")
	w.WriteHeader(http.StatusOK)
}
