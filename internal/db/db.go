package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ArtemNehoda/golang-hello-world/internal/ports"
	_ "github.com/go-sql-driver/mysql"
)

// Good only for this hello world. In a real application, I would obviously use migrations.
// For example, with https://github.com/golang-migrate/migrate
const createTableSQL = `
CREATE TABLE IF NOT EXISTS messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    content VARCHAR(255) NOT NULL,
    author VARCHAR(100) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`

// InitDB opens a MySQL connection, waits for it to be ready (up to 30 retries),
// and creates the messages table if it does not exist.
func InitDB(dsn string, log ports.Logger) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Retry ping up to 30 times with 2s delay to handle container startup timing.
	const maxRetries = 30
	for i := 1; i <= maxRetries; i++ {
		if err = db.Ping(); err == nil {
			break
		}
		log.Printf("Database not ready (attempt %d/%d): %v", i, maxRetries, err)
		if i == maxRetries {
			log.Fatalf("Could not connect to database after %d attempts: %v", maxRetries, err)
		}
		time.Sleep(2 * time.Second)
	}

	log.Println("Database connection established.")

	if _, err = db.Exec(createTableSQL); err != nil {
		log.Fatalf("Failed to create messages table: %v", err)
	}

	log.Println("Messages table ready.")
	return db
}

// CleanData deletes all records from the messages table.
// Useful for tests that need a clean slate.
func CleanData(db *sql.DB, log ports.Logger) error {
	if _, err := db.Exec("DELETE FROM messages"); err != nil {
		return fmt.Errorf("failed to clean messages table: %w", err)
	}
	log.Println("Messages table cleaned.")
	return nil
}

// SeedData inserts sample messages if the messages table is empty.
func SeedData(db *sql.DB, log ports.Logger) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM messages").Scan(&count); err != nil {
		return fmt.Errorf("failed to count messages: %w", err)
	}

	if count > 0 {
		log.Printf("Database already has %d message(s), skipping seed.", count)
		return nil
	}

	seeds := []struct {
		content string
		author  string
	}{
		{"Hello World!", "Artem Nehoda"},
		{"Welcome to Go!", "Artemisio Maltempo"},
		{"I can do it better!", "Artemon Stormbringer"},
	}

	for _, s := range seeds {
		if _, err := db.Exec(
			"INSERT INTO messages (content, author) VALUES (?, ?)",
			s.content, s.author,
		); err != nil {
			return fmt.Errorf("failed to seed message %q: %w", s.content, err)
		}
	}

	log.Printf("Seeded %d messages.", len(seeds))
	return nil
}
