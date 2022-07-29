package gin_test

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	metrics "github.com/yangshun2005/go-prometheus-http/metrics/prometheus"
	"github.com/yangshun2005/go-prometheus-http/middleware"
	ginmiddleware "github.com/yangshun2005/go-prometheus-http/middleware/gin"
)

// GinMiddleware shows how you would create a default middleware factory and use it
// to create a Gin compatible middleware.
func Example_ginMiddleware() {
	// Create our middleware factory with the default settings.
	mdlw := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	// Create our gin instance.
	engine := gin.New()

	// Add our handler and middleware
	h := func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world")
	}
	engine.GET("/", ginmiddleware.Handler("", mdlw), h)

	// Serve metrics from the default prometheus registry.
	log.Printf("serving metrics at: %s", ":8081")
	go func() {
		_ = http.ListenAndServe(":8081", promhttp.Handler())
	}()

	// Serve our handler.
	log.Printf("listening at: %s", ":8080")
	if err := http.ListenAndServe(":8080", engine); err != nil {
		log.Panicf("error while serving: %s", err)
	}
}
