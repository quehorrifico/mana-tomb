package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Connect initializes the global DB connection to mana_tomb.
func Connect() {
	// Step 1: Connect to default "postgres" DB
	postgresDsn := "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"
	tempDB, err := sql.Open("postgres", postgresDsn)
	if err != nil {
		log.Fatal("❌ Could not open default postgres DB:", err)
	}
	defer tempDB.Close()

	// Step 2: Try to create the mana_tomb database
	_, err = tempDB.Exec(`CREATE DATABASE mana_tomb;`)
	if err != nil && err.Error() != `pq: database "mana_tomb" already exists` {
		log.Fatal("❌ Could not create mana_tomb database:", err)
	}

	// Step 3: Connect to the actual mana_tomb DB
	manaTombDsn := "postgres://postgres:password@localhost:5432/mana_tomb?sslmode=disable"
	DB, err = sql.Open("postgres", manaTombDsn)
	if err != nil {
		log.Fatal("❌ Could not connect to mana_tomb DB:", err)
	}

	// Step 4: Verify connection
	if err := DB.Ping(); err != nil {
		log.Fatal("❌ Could not ping mana_tomb DB:", err)
	}

	log.Println("✅ Connected to mana_tomb database")
}

// GetDB returns the active database connection
func GetDB() *sql.DB {
	return DB
}
