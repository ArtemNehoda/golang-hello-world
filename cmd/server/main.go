package main

import (
	"fmt"
	"net/http"

	"github.com/ArtemNehoda/golang-hello-world/internal/config"
	"github.com/ArtemNehoda/golang-hello-world/internal/db"
	"github.com/ArtemNehoda/golang-hello-world/internal/graphql"
	"github.com/ArtemNehoda/golang-hello-world/internal/repositories/mysql"
	"github.com/ArtemNehoda/golang-hello-world/internal/services"
	"github.com/ArtemNehoda/golang-hello-world/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	log := logger.New()

	// Load .env file; ignore error if not present (env vars may be set directly).
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables.")
	}

	// Load configuration from environment variables.
	cfg := config.Load()

	// Initialize database connection and ensure messages table exists.
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	sqlDB := db.InitDB(dsn, log)
	defer sqlDB.Close()

	if err := db.SeedData(sqlDB, log); err != nil {
		log.Fatalf("Failed to seed data: %v", err)
	}

	// Wire up GraphQL resolver with service dependency.
	repo := mysql.NewMessageRepository(sqlDB, log)
	messageService := services.NewMessageService(repo)
	resolver := &graphql.Resolver{
		Service: messageService,
		Logger:  log,
	}

	mux := http.NewServeMux()
	mux.Handle("/graphql", graphql.NewGraphQLHandler(resolver))

	addr := ":" + cfg.ServerPort
	log.Printf("GraphQL server starting on %s", addr)
	log.Printf("GraphQL endpoint:  http://localhost%s/graphql", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
