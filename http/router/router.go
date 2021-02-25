package router

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type ErrorHandlerFunc func(http.ResponseWriter, *http.Request, int)

var defaultErrorHandlerFunc = func(rw http.ResponseWriter, _ *http.Request, i int) {
	http.Error(rw, http.StatusText(i), i)
}

type Router struct {
	re           *regexp.Regexp
	paths        strings.Builder
	ErrorHandler ErrorHandlerFunc
	handlers     map[string]http.Handler
	counter      int
}

func New() *Router {
	return &Router{handlers: make(map[string]http.Handler)}
}

func (r *Router) Route(method, path string, handler http.Handler) {
	if r.paths.Len() > 0 {
		r.paths.WriteByte('|')
	}
	fmt.Fprintf(&r.paths, "(?P<%s%d>%s)", method, r.counter, path)
	r.handlers[strconv.Itoa(r.counter)] = handler
	r.counter++
}

func (r *Router) Compile() http.Handler {
	r.re = regexp.MustCompile("^(?:" + r.paths.String() + ")$")
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		statusCode := http.StatusNotFound
		match := r.re.FindStringSubmatch(req.URL.Path)
		if len(match) == 0 {
			goto theEnd
		}
		for i, name := range r.re.SubexpNames() {
			if len(name) > 0 && len(match[i]) > 0 {
				// if name != "" && match[i] != "" {
				if strings.HasPrefix(name, req.Method) {
					r.handlers[name].ServeHTTP(res, req)
					return
				}
				statusCode = http.StatusMethodNotAllowed
				break
			}
		}
	theEnd:
		if r.ErrorHandler != nil {
			r.ErrorHandler(res, req, statusCode)
			return
		}
		defaultErrorHandlerFunc(res, req, statusCode)
	})
}
