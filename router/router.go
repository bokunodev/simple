package router

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Router struct {
	re              *regexp.Regexp
	paths           strings.Builder
	NotFoundHandler http.HandlerFunc
	handlers        map[string]http.Handler
	counter         int
}

func New() *Router {
	return &Router{handlers: make(map[string]http.Handler)}
}

func (r *Router) Route(path string, handler http.Handler) {
	if r.paths.Len() > 0 {
		r.paths.Write([]byte{'|'})
	}
	fmt.Fprintf(&r.paths, "(?P<%d>%s)", r.counter, path)
	r.handlers[strconv.Itoa(r.counter)] = handler
	r.counter++
}

func (r *Router) Compile() http.Handler {
	r.re = regexp.MustCompile("^(?:" + r.paths.String() + ")$")
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		match := r.re.FindStringSubmatch(req.URL.Path)
		if len(match) == 0 {
			goto theEnd
		}
		for i, name := range r.re.SubexpNames() {
			if name != "" && match[i] != "" {
				r.handlers[name].ServeHTTP(res, req)
				return
			}
		}
	theEnd:
		if r.NotFoundHandler != nil {
			r.NotFoundHandler(res, req)
			return
		}
		http.NotFound(res, req)
	})
}
