package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	contentType = "application/json"
)

// HTTPServer represents an HTTP server.
type HTTPServer struct {
	listenHost string
	debug      bool
	server     *http.Server
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
		log.Debug().Msg("Debug mode enabled")
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()

	healthRoutes := r.Group("/health")
	healthRoutes.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	healthRoutes.GET("/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ready",
		})
	})

	s.server = &http.Server{
		Addr:    s.listenHost,
		Handler: r,
	}

	log.Info().Msgf("HTTP Server started on http://%v\n", s.listenHost)

	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server.
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down HTTP server...")
	return s.server.Shutdown(ctx)
}
