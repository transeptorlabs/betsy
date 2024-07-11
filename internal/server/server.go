package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/transeptorlabs/betsy/internal/mempool"
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
	mempool    *mempool.UserOpMempool
}

// NewHTTPServer creates a new HTTP server.
func NewHTTPServer(listenHost string, debug bool, wallet *wallet.Wallet, mempool *mempool.UserOpMempool) *HTTPServer {
	return &HTTPServer{
		listenHost: listenHost,
		debug:      debug,
		wallet:     wallet,
		mempool:    mempool,
	}
}

// Run starts the HTTP server.
func (s *HTTPServer) Run() error {
	log.Info().Msg("Starting HTTP server...")

	if s.debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.LoadHTMLGlob("ui/templates/*")
	router.Static("assets/css", "ui/assets/css")
	router.Static("assets/js", "ui/assets/js")
	router.Static("assets/img", "ui/assets/img")

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

	// Dashboard
	router.GET("/dashboard", func(c *gin.Context) {
		c.HTML(http.StatusOK, "base", gin.H{
			"title": "Dashboard",
		})
	})

	router.GET("/accounts", func(c *gin.Context) {
		accounts, err := s.wallet.GetDevAccounts(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}

		c.HTML(http.StatusOK, "accounts", gin.H{
			"accounts": accounts,
		})
	})

	router.GET("/mempool", func(c *gin.Context) {
		c.HTML(http.StatusOK, "mempool", gin.H{
			"userOps": s.mempool.GetUserOps(),
		})
	})

	router.GET("/bundles", func(c *gin.Context) {
		c.HTML(http.StatusOK, "bundles", gin.H{
			"bundles": nil,
		})
	})

	router.NoRoute(func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/dashboard")
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
