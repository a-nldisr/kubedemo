package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var someRandomMetric = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "some_random_metric", Help: "A random static metric in your Kubernetes to play with"})

type application struct {
	readinessFailure bool
	livenessFailure  bool
	infoLog          *log.Logger
	errorLog         *log.Logger
	debugLog         *log.Logger
}

func (app *application) randomChannel(c chan bool) {
	app.infoLog.Println("Starting time based random metric generator")

	// Time in seconds
	pollInterval := 10

	timerCh := time.Tick(time.Duration(pollInterval) * time.Second)
	// Time based loop to generate Global variable
	for range timerCh {
		select {
		// when shutdown is received we break
		case <-c:
			app.infoLog.Println("Received shutdown, stopping timer")
			break
		default:
			app.debugLog.Println("Generating number... ")
			min := 0
			max := 10
			num := float64(rand.Intn(max-min) + min)
			s := fmt.Sprintf("%f", num)
			app.debugLog.Println("Generated random number: " + string(s))

			someRandomMetric.Set(num)
		}
	}
}

func init() {
	prometheus.MustRegister(someRandomMetric)
	// Set the metric to 0 at startup
	someRandomMetric.Set(0)
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	debugLog := log.New(os.Stdout, "DEBUG\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		readinessFailure: false,
		livenessFailure:  false,
		infoLog:          infoLog,
		errorLog:         errorLog,
		debugLog:         debugLog,
	}

	app.infoLog.Println("Version 0.0.4")

	port := ":8090"

	server := &http.Server{
		Addr:         port,
		ErrorLog:     app.errorLog,
		Handler:      app.routes(),
		IdleTimeout:  15 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.infoLog.Printf("Starting server on %v\n", port)

	// Creating the channel in a go routine to ensure the timer runs concurrently in the background
	notificationChannel := make(chan bool)
	go app.randomChannel(notificationChannel)

	// Setting clean shutdown mechanisms
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		// quit blocks till the process gets a shutdown signal
		<-quit
		// Send notification to close the time loop in childChan
		notificationChannel <- true
		// Notifying the server shutdown is received.
		app.infoLog.Println("Server is shutting down.")

		// Setting up context group
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			app.errorLog.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}

		// Closing channel
		close(done)
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		app.errorLog.Fatalf("Could not listen on %s: %v\n", port, err)
	} // Block till server is closed, send notification after
	<-done
	app.infoLog.Println("Server stopped")

}
