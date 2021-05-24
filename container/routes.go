package main

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// routes returns a router with all paths.
func (app *application) routes() *mux.Router {
	r := mux.NewRouter()
	r.Use(app.printFunctionName)
	r.HandleFunc("/", app.hello).Name("hello")
	r.HandleFunc("/ip", app.ip).Name("ip")
	r.HandleFunc("/headers", app.headers).Name("headers")
	r.HandleFunc("/livez", app.livez).Name("livez")
	r.HandleFunc("/readyz", app.readyz).Name("readyz")
	r.HandleFunc("/config/crash", app.crash).Name("crash")
	r.HandleFunc("/config/livezfailure", app.livezFailure).Name("livezFailure")
	r.HandleFunc("/config/readyzfailure", app.readyzFailure).Name("readyzFailure")
	r.HandleFunc("/factorial", app.factorial).Name("factorial")
	r.Handle("/metrics", promhttp.Handler()).Name("metrics")
	return r
}
