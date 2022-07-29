package fasthttp_test

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	metrics "github.com/yangshun2005/go-prometheus-http/metrics/prometheus"
	"github.com/yangshun2005/go-prometheus-http/middleware"
	fasthttpMiddleware "github.com/yangshun2005/go-prometheus-http/middleware/fasthttp"
)

// FasthttpMiddleware shows how you would create a default middleware
// factory and use it to create a fasthttp compatible middleware.
func Example_fasthttpMiddleware() {
	// Create our middleware factory with the default settings.
	mdlw := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	// Add our handler and middleware
	h := func(rCtx *fasthttp.RequestCtx) {
		rCtx.SetStatusCode(fasthttp.StatusOK)
		rCtx.SetBodyString("OK")
	}

	// Create our fasthttp instance.
	srv := &fasthttp.Server{
		Handler: fasthttpMiddleware.Handler("", mdlw, h),
	}

	// Serve metrics from the default prometheus registry.
	log.Printf("serving metrics at: %s", ":8081")
	go func() {
		_ = http.ListenAndServe(":8081", promhttp.Handler())
	}()

	// Serve our handler.
	log.Printf("listening at: %s", ":8080")
	if err := srv.ListenAndServe(":8080"); err != nil {
		log.Panicf("error while serving: %s", err)
	}
}
