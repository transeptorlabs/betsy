package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	contentType = "application/json"
)

// HTTPServer represents an HTTP server.
type HTTPServer struct {
	listenHost string
	debug      bool
}

// NewHTTPServer creates a new HTTP server.
func NewHTTPServer(listenHost string, debug bool) *HTTPServer {
	return &HTTPServer{
		listenHost: listenHost,
		debug:      debug,
	}
}

// Run starts the HTTP server.
func (s *HTTPServer) Run() error {
	if s.debug {
		fmt.Println("Debug mode enabled")
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()

	healthRoutes := r.Group("/health")
	healthRoutes.GET("/", func(c *gin.Context) {
		fmt.Println("Health check")
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	healthRoutes.GET("/ready", func(c *gin.Context) {
		fmt.Println("Readiness check")
		c.JSON(200, gin.H{
			"status": "ready",
		})
	})

	fmt.Printf("HTTP Server started on http://%v\n", s.listenHost)

	return r.Run(s.listenHost)
}
