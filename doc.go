/*
Package gohttpmetrics knows how to measure http metrics in different metric formats,
it comes with a middleware that can be used for different frameworks and also the
the main Go net/http handler:
	package main

	import (
		"log"
		"net/http"

		"github.com/prometheus/client_golang/prometheus/promhttp"
		httpmetrics "github.com/yangshun2005/go-prometheus-http/metrics/prometheus"
		httpmiddleware "github.com/yangshun2005/go-prometheus-http/middleware"
		httpstdmiddleware "github.com/yangshun2005/go-prometheus-http/middleware/std"
	)

	func main() {
		// Create our middleware.
		mdlw := httpmiddleware.New(httpmiddleware.Config{
			Recorder: httpmetrics.NewRecorder(httpmetrics.Config{}),
		})

		// Our handler.
		myHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("hello world!"))
		})
		h := httpstdmiddleware.Handler("", mdlw, myHandler)

		// Serve metrics.
		log.Printf("serving metrics at: %s", ":9090")
		go http.ListenAndServe(":9090", promhttp.Handler())

		// Serve our handler.
		log.Printf("listening at: %s", ":8080")
		if err := http.ListenAndServe(":8080", h); err != nil {
			log.Panicf("error while serving: %s", err)
		}
	}
*/
package gohttpmetrics

// blank imports help docs.
import (
	// Import metrics package.
	_ "github.com/yangshun2005/go-prometheus-http/metrics"
	// Import middleware package.
	_ "github.com/yangshun2005/go-prometheus-http/middleware"
)
