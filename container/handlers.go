package main

import (
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strconv"
)

func (app *application) hello(w http.ResponseWriter, req *http.Request) {
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

	var msg string

	if app.livenessFailure {
		msg = "enabled failure into /livez endpoint"
	} else {
		msg = "disabled failure mode for /livez endpoint"
	}

	app.debugLog.Println(msg)
	fmt.Fprintln(w, msg)
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

	var msg string

	if app.readinessFailure {
		msg = "enabled failure into /readyz endpoint"
	} else {
		msg = "disabled failure mode for /readyz endpoint"
	}

	app.debugLog.Println(msg)
	fmt.Fprintln(w, msg)
}

func (app *application) headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func (app *application) factorial(w http.ResponseWriter, req *http.Request) {
	var num int64 = 50000

	s := req.FormValue("number")
	if s != "" {
		n, err := strconv.Atoi(s)
		if err == nil && n > 0 {
			num = int64(n)
		} else {
			status := http.StatusBadRequest
			w.WriteHeader(status)
			fmt.Fprintln(w, status, http.StatusText(status))
			fmt.Fprintf(w, "Only numbers > 0 are accepted")
			return
		}
	}

	_ = app.getFactorial(big.NewInt(num))
	fmt.Fprintf(w, "Calculated factorial for %d\n", num)
}
