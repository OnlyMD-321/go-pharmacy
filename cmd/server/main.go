package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/OnlyMD-321/go-pharmacy/internal/api"
	"github.com/OnlyMD-321/go-pharmacy/internal/config"
	"github.com/OnlyMD-321/go-pharmacy/internal/db"
	"github.com/OnlyMD-321/go-pharmacy/internal/firebase"
	"github.com/OnlyMD-321/go-pharmacy/internal/middlewares"
)

func main() {
	// Load environment variables
	config.Load()

	// Initialize PostgreSQL connection pool
	db.InitDB()
	defer db.CloseDB()

	// Initialize Firebase App
	firebase.InitFirebase()

	// Setup Gin router
	router := gin.Default()

	// Public health check route
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Initialize handlers with DB pool
	userHandler := api.NewUserHandler(db.Pool)
	inventoryHandler := api.NewInventoryHandler(db.Pool)
	saleHandler := api.NewSaleHandler(db.Pool)

	// API routes group with Firebase authentication middleware
	apiGroup := router.Group("/api")
	apiGroup.Use(middlewares.FirebaseAuthMiddleware())

	// User profile
	apiGroup.GET("/profile", userHandler.GetProfile)

	// Inventory routes (admin, pharmacist)
	apiGroup.GET("/inventory", middlewares.NewRBACMiddleware(db.Pool, "admin", "pharmacist"), inventoryHandler.GetInventory)
	apiGroup.POST("/inventory", middlewares.NewRBACMiddleware(db.Pool, "admin"), inventoryHandler.CreateInventory)

	// Sales routes (admin, pharmacist, seller)
	apiGroup.GET("/sales", middlewares.NewRBACMiddleware(db.Pool, "admin", "pharmacist"), saleHandler.GetSales)
	apiGroup.POST("/sales", middlewares.NewRBACMiddleware(db.Pool, "seller"), saleHandler.CreateSale)

	// Start server with graceful shutdown
	srv := &http.Server{
		Addr:    ":" + config.AppConfig.Port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("🔥 listen: %s\n", err)
		}
	}()
	log.Printf("🚀 Server running on port %s", config.AppConfig.Port)

	// Wait for termination signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("⚙️ Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
