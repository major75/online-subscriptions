package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/major75/online-subscriptions/database"
	_ "github.com/major75/online-subscriptions/docs"
	"github.com/major75/online-subscriptions/internal"
	"github.com/major75/online-subscriptions/pkg/logger"
)

// @title Subscriptions Service
// @version 1.0.0
// @BasePath /
// @description A production-ready Go subscription service
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @tag.name system

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Failed to load .env file: %v. Will use environment variables\n", err)
	}

	var exist bool
	serviceName, exist := os.LookupEnv("SERVICE_NAME")
	if !exist {
		fmt.Println("Env variable SERVICE_NAME is not set")
		os.Exit(1)
	}

	logLevel := os.Getenv("LOG_LEVEL")
	log, err := logger.New(logLevel, serviceName)
	if err != nil {
		fmt.Printf("failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	serverHost, exist := os.LookupEnv("SERVER_HOST")
	if !exist {
		log.Fatal("Env variable SERVER_HOST is not set")
	}
	serverPort, exist := os.LookupEnv("SERVER_PORT")
	if !exist {
		log.Fatal("Env variable SERVER_PORT is not set")
	}

	dsn, exist := os.LookupEnv("DATABASE_URL")
	if !exist {
		log.Fatal("Env variable DATABASE_URL is not set")
	}

	maxConnections, err := strconv.Atoi(os.Getenv("DATABASE_MAX_CONNECTIONS"))
	if err != nil {
		log.Warn("Failed to parse database max connections from config. Use default value")
		maxConnections = 20
	}
	connectionLifetime, err := strconv.Atoi(os.Getenv("DATABASE_CONNECTION_LIFETIME"))
	if err != nil {
		log.Warn("Failed to parse database connection lifetime from config. Use default value")
		connectionLifetime = 10
	}

	dbCfg := database.DBCfg{
		Dsn:                dsn,
		MaxConnections:     int32(maxConnections),
		ConnectionLifetime: time.Duration(connectionLifetime) * time.Minute,
	}

	db, err := database.NewPostgreDB(context.Background(), dbCfg, log)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}

	log.Info("Starting service")

	err = database.RunMigrations(db, log)
	if err != nil {
		log.Fatal("Failed to apply SQL migrations", "error", err)
	}

	r := router.NewRouter(log, db)

	addr := fmt.Sprintf("%s:%s", serverHost, serverPort)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start server in go routine
	go func() {
		log.Info("Server listening", "address", addr)
		err = srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Server failed to start", "error", err)
		}
	}()

	// Graceful shutdown
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGTERM)
	<-quitCh

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server stopped forcibly by timeout", "error", err)
	}

	log.Info("Server stopped gracefully")
}
