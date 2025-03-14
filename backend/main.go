package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/quehorrifico/mana-tomb/handlers"
	"github.com/quehorrifico/mana-tomb/utils"
)

var (
	db     *sql.DB
	db_err error
)

// Connects to PostgreSQL (creates database if missing)
func connectDB() {
	postgresDsn := "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"

	// Connect to the default "postgres" database to check for "mana_tomb"
	db, db_err = sql.Open("postgres", postgresDsn)
	if db_err != nil {
		// return nil, fmt.Errorf("error connecting to PostgreSQL: %w", err)
		log.Fatal("❌ Error connecting to PostgreSQL:", db_err)
	}
	defer db.Close()

	// Ensure the "mana_tomb" database exists
	_, err := db.Exec("CREATE DATABASE mana_tomb;")
	if err != nil && err.Error() != `pq: database "mana_tomb" already exists` {
		// return nil, fmt.Errorf("error creating database: %w", err)
		log.Fatal("❌ Error creating database:", err)
	}

	// Now connect to the mana_tomb database
	manaTombDsn := "postgres://postgres:password@localhost:5432/mana_tomb?sslmode=disable"
	db, err = sql.Open("postgres", manaTombDsn)
	if err != nil {
		// return nil, fmt.Errorf("error connecting to mana_tomb: %w", err)
		log.Fatal("❌ Error connecting to mana_tomb:", err)
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		// return nil, fmt.Errorf("error verifying database connection: %w", err)
		log.Fatal("❌ Error verifying database connection:", err)
	}

	log.Println("✅ Connected to PostgreSQL database!")
	// return db, nil
}

// Ensure required tables exist
func ensureTables(db *sql.DB) {
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

	_, err := db.Exec(createBulkDataTable)
	if err != nil {
		log.Fatalf("❌ Error creating bulk_data table: %v", err)
	} else {
		log.Println("✅ bulk_data table checked/created successfully!")
	}
}

// Get a random card from the database
func randomCard(w http.ResponseWriter, r *http.Request) {
	enableCORS(w) // Add CORS headers
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Use the existing db connection
	query := `
		SELECT name, mana_cost, oracle_text
		FROM oracle_cards
		ORDER BY random()
		LIMIT 1;
	`
	var card handlers.OracleCard

	// Fetch a random card
	err := db.QueryRow(query).Scan(&card.Name, &card.ManaCost, &card.OracleText)
	if err != nil {
		http.Error(w, "Error fetching card", http.StatusInternalServerError)
		log.Println("❌ Error fetching random card:", err)
		return
	}

	// Convert to JSON and return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(card)
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func main() {
	// Connect to the database
	// db, err := connectDB()
	// if err != nil {
	// 	log.Fatalf("❌ Database connection failed: %v", err)
	// }
	// defer db.Close()
	connectDB()

	// Ensure tables exist
	ensureTables(db)

	// Start the bulk data fetch scheduler
	utils.StartScheduler(db)

	// API endpoint
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Mana Tomb API is running!"))
	})
	http.HandleFunc("/random-card", randomCard)

	log.Println("🚀 Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
