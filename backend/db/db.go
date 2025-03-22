// package db

// import (
// 	"database/sql"
// 	"log"

// 	_ "github.com/lib/pq"
// )

// var DB *sql.DB

// // Connect to PostgreSQL, create mana_tomb DB if missing
// func Connect() error {
// 	postgresDsn := "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"

// 	// Connect to the default "postgres" database to check for "mana_tomb"
// 	DB, err := sql.Open("postgres", postgresDsn)
// 	if err != nil {
// 		return err
// 	}
// 	defer DB.Close()

// 	// Ensure the "mana_tomb" database exists
// 	_, err = DB.Exec("CREATE DATABASE mana_tomb;")
// 	if err != nil && err.Error() != `pq: database "mana_tomb" already exists` {
// 		return err
// 	}

// 	// Now connect to the mana_tomb database
// 	manaTombDsn := "postgres://postgres:password@localhost:5432/mana_tomb?sslmode=disable"
// 	DB, err = sql.Open("postgres", manaTombDsn)
// 	if err != nil {
// 		return err
// 	}

// 	// Verify the connection
// 	if err := DB.Ping(); err != nil {
// 		return err
// 	}

// 	log.Println("✅ Connected to PostgreSQL database!")
// 	return nil
// }

// func EnsureTables() error {
// 	createBulkDataTable := `
// 	CREATE TABLE IF NOT EXISTS bulk_data (
// 	    id UUID PRIMARY KEY,
// 	    uri TEXT NOT NULL,
// 	    type TEXT NOT NULL,
// 	    name TEXT NOT NULL,
// 	    description TEXT NOT NULL,
// 	    download_uri TEXT NOT NULL,
// 	    updated_at TIMESTAMP NOT NULL,
// 	    size BIGINT NOT NULL,
// 	    content_type TEXT NOT NULL,
// 	    content_encoding TEXT NOT NULL
// 	);`

// 	_, err := DB.Exec(createBulkDataTable)
// 	if err != nil {
// 		return err
// 	}

// 	// Log success
// 	log.Println("✅ bulk_data table checked/created successfully!")
// 	return nil
// }

// // GetDB returns the shared DB instance
// func GetDB() *sql.DB {
// 	return DB
// }

package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// DB is the final, persistent connection to mana_tomb
var DB *sql.DB

// Connect ensures mana_tomb DB exists, then opens the final DB connection
func Connect() {
	// 1) Connect to the default "postgres" DB
	postgresDsn := "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"
	tempDB, err := sql.Open("postgres", postgresDsn)
	if err != nil {
		log.Fatal("❌ Error opening default postgres DB:", err)
	}

	// Attempt to create mana_tomb DB if it doesn't exist
	_, err = tempDB.Exec(`CREATE DATABASE mana_tomb;`)
	// If error is anything but "already exists", fail
	if err != nil && err.Error() != `pq: database "mana_tomb" already exists` {
		log.Fatalf("❌ Error creating mana_tomb DB: %v", err)
	}

	// Close the ephemeral connection
	tempDB.Close()

	// 2) Now connect to the actual mana_tomb DB
	manaTombDsn := "postgres://postgres:password@localhost:5432/mana_tomb?sslmode=disable"
	DB, err = sql.Open("postgres", manaTombDsn)
	if err != nil {
		log.Fatal("❌ Error connecting to mana_tomb:", err)
	}

	if pingErr := DB.Ping(); pingErr != nil {
		log.Fatal("❌ Error verifying mana_tomb connection:", pingErr)
	}

	log.Println("✅ Connected to the mana_tomb database!")
}

// EnsureTables runs any CREATE TABLE logic you need
func EnsureTables() {
	createBulkDataTable := `
	CREATE TABLE IF NOT EXISTS bulk_data (
	    id UUID PRIMARY KEY,
	    uri TEXT NOT NULL,
	    type TEXT NOT NULL,
	    name TEXT NOT NULL,
	    description TEXT NOT NULL,
	    download_uri TEXT NOT NULL,
	    updated_at TIMESTAMP NOT NULL,
	    size BIGINT NOT NULL,
	    content_type TEXT NOT NULL,
	    content_encoding TEXT NOT NULL
	);`
	_, err := DB.Exec(createBulkDataTable)
	if err != nil {
		log.Fatalf("❌ Error creating bulk_data table: %v", err)
	} else {
		log.Println("✅ bulk_data table checked/created successfully!")
	}

	// Add any other tables (e.g., orchard_cards, oracle_cards, etc.)
}

// GetDB returns the persistent DB
func GetDB() *sql.DB {
	return DB
}
