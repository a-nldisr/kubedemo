package main

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func (app *application) printFunctionName(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		remoteIP := req.Header.Get("X-Forwarded-For")
		if remoteIP == "" {
			remoteIP = strings.Split(req.RemoteAddr, ":")[0]
		}

		routeName := mux.CurrentRoute(req).GetName()
		app.debugLog.Printf("%s called %q", remoteIP, routeName)
		next.ServeHTTP(w, req)
	})
}
