package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var (
	ErrBadRequest = errors.New("Bad request")
	ErrForbidden  = errors.New("Forbidden")
	ErrNotFound   = errors.New("Not found")
	ErrUnknown    = errors.New("Unknown error")
)

type Router struct {
	entries []entry

	debug bool
	delay int
}

type Handlers = map[string]func(*http.Request) (interface{}, error)

type entry struct {
	path     string
	exact    bool
	handlers Handlers
}

func NewRouter(debug bool, delay int) Router {
	return Router{[]entry{}, debug, delay}
}

func (router *Router) Handle(path string, handlers Handlers) {
	router.entries = append(router.entries, entry{path, !strings.HasSuffix(path, "/"), handlers})
}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if router.debug {
		fmt.Printf("%s %s -> ", r.Method, r.URL.Path)
	}

	for _, e := range router.entries {
		if (e.exact && r.URL.Path == e.path) || (!e.exact && strings.HasPrefix(r.URL.Path, e.path)) {
			handler := e.handlers[r.Method]
			if handler != nil {
				obj, err := handler(r)
				if obj != nil && err == nil {
					router.writeJson(w, obj)
				} else {
					router.writeError(w, err)
				}
				return
			}

			if r.Method == http.MethodOptions {
				router.writeOptions(w, e)
				return
			}

			router.writeError(w, ErrBadRequest)
			return
		}
	}

	router.writeError(w, ErrNotFound)
}

func (router *Router) writeOptions(w http.ResponseWriter, e entry) {
	methods := ""
	for method := range e.handlers {
		methods += method + ", "
	}
	methods += http.MethodOptions

	if router.debug {
		fmt.Println("OK")
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", methods)
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	w.Header().Set("Access-Control-Max-Age", "86400")
	w.WriteHeader(http.StatusOK)
}

func (router *Router) writeJson(w http.ResponseWriter, obj interface{}) {
	b, err := json.Marshal(obj)
	if err != nil {
		router.writeError(w, err)
		return
	}

	if router.delay > 0 {
		time.Sleep(time.Duration(router.delay) * time.Millisecond)
	}

	if router.debug {
		fmt.Println(string(b))
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (router *Router) writeError(w http.ResponseWriter, err error) {
	if err == nil {
		err = ErrUnknown
	}

	var status int
	switch err {
	case ErrBadRequest:
		status = http.StatusBadRequest
	case ErrForbidden:
		status = http.StatusForbidden
	case ErrNotFound:
		status = http.StatusNotFound
	default:
		status = http.StatusInternalServerError
	}

	if router.delay > 0 {
		time.Sleep(time.Duration(router.delay) * time.Millisecond)
	}

	if router.debug {
		fmt.Printf("%d %s\n", status, err.Error())
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}
