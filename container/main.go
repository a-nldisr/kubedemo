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
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var someRandomMetric = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "some_random_metric", Help: "A random static metric in your Kubernetes to play with"})

type application struct {
	readinessFailure bool
	livenessFailure  bool
}

func randomChannel(c chan bool) {
	fmt.Println("Starting time based random metric generator")

	// Time in seconds
	pollInterval := 10

	timerCh := time.Tick(time.Duration(pollInterval) * time.Second)
	// Time based loop to generate Global variable
	for range timerCh {
		select {
		// when shutdown is received we break
		case <-c:
			fmt.Println("Received shutdown, stopping timer")
			break
		default:
			fmt.Println("Generating number... ")
			min := 0
			max := 10
			num := float64(rand.Intn(max-min) + min)
			s := fmt.Sprintf("%f", num)
			fmt.Println("Generated random number: " + string(s))

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

	fmt.Println("Version 0.0.4")

	port := ":8090"

	server := &http.Server{
		Addr:         port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	app := &application{}

	fmt.Printf("Starting server on %v\n", port)

	// Creating the channel in a go routine to ensure the timer runs concurrently in the background
	notificationChannel := make(chan bool)
	go randomChannel(notificationChannel)

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
		fmt.Println("Server is shutting down.")

		// Setting up context group
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}

		// Closing channel
		close(done)
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", app.hello)
	http.HandleFunc("/headers", app.headers)
	http.HandleFunc("/livez", app.livez)
	http.HandleFunc("/readyz", app.readyz)
	http.HandleFunc("/config/livezfailure", app.livezFailure)
	http.HandleFunc("/config/readyzfailure", app.readyzFailure)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", port, err)
	} // Block till server is closed, send notification after
	<-done
	fmt.Println("Server stopped")

}
