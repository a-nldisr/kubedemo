package main

import (
	"fmt"
	"net/http"
	"os"
)

func (app *application) hello(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Hello function")

	fmt.Fprintf(w, "Hi Vandebron\n")
	fmt.Fprintf(w, "Our sealed secret is: "+os.Getenv("SECRET"))
	fmt.Fprintf(w, "\nOur enviroment variable is: "+os.Getenv("FOO"))
}

func (app *application) livez(w http.ResponseWriter, req *http.Request) {
	status := http.StatusOK

	if app.livenessFailure {
		status = http.StatusBadGateway
	}

	w.WriteHeader(status)
	fmt.Fprintln(w, status, http.StatusText(status))
}

func (app *application) livezFailure(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(405)
		fmt.Fprintln(w, "Method Not Allowed")
	}
	app.livenessFailure = !app.livenessFailure
}

func (app *application) readyz(w http.ResponseWriter, req *http.Request) {
	status := http.StatusOK

	if app.readinessFailure {
		status = http.StatusTooManyRequests
	}

	w.WriteHeader(status)
	fmt.Fprintln(w, status, http.StatusText(status))
}

func (app *application) readyzFailure(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(405)
		fmt.Fprintln(w, "Method Not Allowed")
	}
	app.readinessFailure = !app.readinessFailure
}

func (app *application) headers(w http.ResponseWriter, req *http.Request) {
	fmt.Println("headers function")

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}
