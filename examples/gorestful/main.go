package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	gorestful "github.com/emicklei/go-restful/v3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/yangshun2005/go-prometheus-http/metrics/prometheus"
	"github.com/yangshun2005/go-prometheus-http/middleware"
	gorestfulmiddleware "github.com/yangshun2005/go-prometheus-http/middleware/gorestful"
)

const (
	srvAddr     = ":8080"
	metricsAddr = ":8081"
)

func main() {
	// Create our middleware.
	mdlw := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	// Create our gorestful instance.
	c := gorestful.NewContainer()

	// Add the middleware for all routes.
	c.Filter(gorestfulmiddleware.Handler("", mdlw))

	// Add our handler.
	ws := &gorestful.WebService{}
	ws.Produces(gorestful.MIME_JSON)

	ws.Route(ws.GET("/").To(func(_ *gorestful.Request, resp *gorestful.Response) {
		resp.WriteEntity("Hello world")
	}))
	ws.Route(ws.GET("/wrong").To(func(_ *gorestful.Request, resp *gorestful.Response) {
		resp.WriteHeaderAndEntity(http.StatusTooManyRequests, "oops")
	}))
	c.Add(ws)

	// Serve our handler.
	go func() {
		log.Printf("server listening at %s", srvAddr)
		if err := http.ListenAndServe(srvAddr, c); err != nil {
			log.Panicf("error while serving: %s", err)
		}
	}()

	// Serve our metrics.
	go func() {
		log.Printf("metrics listening at %s", metricsAddr)
		if err := http.ListenAndServe(metricsAddr, promhttp.Handler()); err != nil {
			log.Panicf("error while serving metrics: %s", err)
		}
	}()

	// Wait until some signal is captured.
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGTERM, syscall.SIGINT)
	<-sigC
}
