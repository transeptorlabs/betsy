package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/transeptorlabs/betsy/wallet"
)

const (
	contentType = "application/json"
)

// HTTPServer represents an HTTP server.
type HTTPServer struct {
	listenHost string
	debug      bool
	server     *http.Server
	wallet     *wallet.Wallet
}

// NewHTTPServer creates a new HTTP server.
func NewHTTPServer(listenHost string, debug bool, wallet *wallet.Wallet) *HTTPServer {
	return &HTTPServer{
		listenHost: listenHost,
		debug:      debug,
		wallet:     wallet,
	}
}

// Run starts the HTTP server.
func (s *HTTPServer) Run() error {
	if s.debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.LoadHTMLGlob("templates/**/*")

	// Health group
	healthRoutes := router.Group("/health")
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

	// Dashboard group
	dashboardRoutes := router.Group("/dashboard")
	dashboardRoutes.GET("/accounts", func(c *gin.Context) {
		accounts, err := s.wallet.GetDevAccounts(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}

		c.HTML(http.StatusOK, "accounts/index.html", gin.H{
			"title":    "Accounts",
			"accounts": accounts,
		})
	})

	dashboardRoutes.GET("/userop-mempool", func(c *gin.Context) {
		c.HTML(http.StatusOK, "mempool/index.html", gin.H{
			"title": "UserOp Mempool",
		})
	})

	dashboardRoutes.GET("/bundles", func(c *gin.Context) {
		c.HTML(http.StatusOK, "bundles/index.html", gin.H{
			"title": "Bundler Bundles",
		})
	})

	s.server = &http.Server{
		Addr:         s.listenHost,
		Handler:      router.Handler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server.
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down HTTP server...")
	return s.server.Shutdown(ctx)
}
