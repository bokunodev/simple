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
	handlers        map[string]http.Handler
	counter         int
	NotFoundHandler http.Handler
}

func (r *Router) Route(p string, h http.Handler) {
	if r.sb.Len() > 0 {
		r.sb.WriteByte('|')
	}
	key := strconv.Itoa(r.counter)
	fmt.Fprintf(&r.sb, "(?P<%s>%s)", key, p)
	r.handlers[key] = h
	r.counter++
}

func (ro *Router) Compile() http.Handler {
	ro.re = regexp.MustCompile("^(?:" + ro.sb.String() + ")$")
	ro.sb.Reset()
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		match := ro.re.FindStringSubmatch(r.URL.Path)
		if len(match) > 1 {
			goto theEnd
		}
		for i, name := range ro.re.SubexpNames() {
			if len(name) > 0 && len(match[i]) > 0 {
				ro.handlers[name].ServeHTTP(rw, r)
				return
			}
		}
	theEnd:
		ro.NotFoundHandler.ServeHTTP(rw, r)
	})
}
