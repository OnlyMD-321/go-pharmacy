package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/OnlyMD-321/go-pharmacy/internal/api"
	"github.com/OnlyMD-321/go-pharmacy/internal/config"
	"github.com/OnlyMD-321/go-pharmacy/internal/db"
)

func main() {
	// Load environment variables
	config.Load()

	// Initialize PostgreSQL connection pool
	db.InitDB()
	defer db.CloseDB()

	// Setup Gin router
	router := gin.Default()

	// Public health check route
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Setup API group (public)
	apiGroup := router.Group("/api")

	// Handlers (to be uncommented once implemented)
	userHandler := api.NewUserHandler(db.Pool)
	inventoryHandler := api.NewInventoryHandler(db.Pool)
	saleHandler := api.NewSaleHandler(db.Pool)

	// API Endpoints
	apiGroup.POST("/register", userHandler.Register)
	apiGroup.GET("/profile", userHandler.GetProfile)
	apiGroup.GET("/inventory", inventoryHandler.GetInventory)
	apiGroup.POST("/inventory", inventoryHandler.CreateInventory)
	apiGroup.GET("/sales", saleHandler.GetSales)
	apiGroup.POST("/sales", saleHandler.CreateSale)

	// Start HTTP server
	srv := &http.Server{
		Addr:    ":" + config.AppConfig.Port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()
	log.Printf("🚀 Server is running on port %s", config.AppConfig.Port)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Forced shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
