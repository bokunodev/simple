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
	sb              strings.Builder
	NotFoundHandler http.HandlerFunc
	handlers        map[string]http.Handler
	counter         int
}

func New() *Router {
	return &Router{
		handlers:        make(map[string]http.Handler),
		NotFoundHandler: http.NotFound,
	}
}

func (r *Router) Route(method, path string, handler http.Handler) {
	if r.sb.Len() > 0 {
		r.sb.WriteByte('|')
	}

	key := strconv.FormatInt(int64(r.counter), 16)
	fmt.Fprintf(&r.sb, "(?P<%s>%s)", key, path)
	r.handlers[key] = handler
	r.counter++
}

func (r *Router) Compile() http.Handler {
	r.re = regexp.MustCompile("^(?:" + r.sb.String() + ")$")
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		match := r.re.FindStringSubmatch(req.URL.Path)
		if len(match) == 0 {
			goto theEnd
		}
		for i, name := range r.re.SubexpNames() {
			if len(name) > 0 && len(match[i]) > 0 {
				r.handlers[name].ServeHTTP(res, req)
				return
			}
		}
	theEnd:
		r.NotFoundHandler(res, req)
	})
}
