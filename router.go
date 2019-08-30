package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

var (
	ErrBadRequest = errors.New("Bad request")
	ErrNotFound   = errors.New("Not found")
	ErrUnknown    = errors.New("Unknown error")
)

type Router struct {
	entries []entry
}

type Handlers = map[string]func(*http.Request) (interface{}, error)

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
				obj, err := handler(r)
				if obj != nil && err == nil {
					writeJson(w, obj)
				} else {
					writeError(w, err)
				}
				return
			}

			if r.Method == http.MethodOptions {
				writeOptions(w, e)
				return
			}

			writeError(w, ErrBadRequest)
			return
		}
	}

	writeError(w, ErrNotFound)
}

func writeOptions(w http.ResponseWriter, e entry) {
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

func writeJson(w http.ResponseWriter, obj interface{}) {
	b, err := json.Marshal(obj)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func writeError(w http.ResponseWriter, err error) {
	if err == nil {
		err = ErrUnknown
	}

	var status int
	switch err {
	case ErrBadRequest:
		status = http.StatusBadRequest
	case ErrNotFound:
		status = http.StatusNotFound
	default:
		status = http.StatusInternalServerError
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}
