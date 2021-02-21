package main

import (
	"net/http"

	"github.com/bokunodev/simple/router"
	"github.com/pkg/profile"
)

func main() {
	defer profile.Start(
		profile.ProfilePath("."),
		profile.MemProfileRate(1),
		profile.MemProfileHeap,
	).Stop()

	r := router.New()
	r.Route(`/user/name/\d{1,2}`, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(r.URL.String()))
	}))
	r.Compile()
	// http.ListenAndServe(":5000", r.Compile())
}
