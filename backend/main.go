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
		log.Fatal("❌ Error connecting to PostgreSQL:", db_err)
	}
	defer db.Close()

	// Ensure the "mana_tomb" database exists
	_, err := db.Exec("CREATE DATABASE mana_tomb;")
	if err != nil && err.Error() != `pq: database "mana_tomb" already exists` {
		log.Fatal("❌ Error creating database:", err)
	}

	// Now connect to the mana_tomb database
	manaTombDsn := "postgres://postgres:password@localhost:5432/mana_tomb?sslmode=disable"
	db, err = sql.Open("postgres", manaTombDsn)
	if err != nil {
		log.Fatal("❌ Error connecting to mana_tomb:", err)
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		log.Fatal("❌ Error verifying database connection:", err)
	}

	log.Println("✅ Connected to PostgreSQL database!")
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
		SELECT name, mana_cost, image_uris, type_line, oracle_text, set, set_name, set_uri, set_id, set_type, set_search_uri, scryfall_set_uri
		FROM oracle_cards
		ORDER BY random()
		LIMIT 1;
	`
	var card handlers.OracleCard
	var imageURIsJSON []byte

	// Fetch a random card
	err := db.QueryRow(query).Scan(&card.Name, &card.ManaCost, &imageURIsJSON, &card.TypeLine, &card.OracleText, &card.Set, &card.SetName, &card.SetURI, &card.SetID, &card.SetType, &card.SetSearchURI, &card.ScryfallSetURI)
	if err != nil {
		http.Error(w, "Error fetching card", http.StatusInternalServerError)
		log.Println("❌ Error fetching random card:", err)
		return
	}

	// Convert JSONB data from database into Go struct
	if err := json.Unmarshal(imageURIsJSON, &card.ImageURIs); err != nil {
		http.Error(w, "Error decoding image URIs", http.StatusInternalServerError)
		log.Println("❌ Error decoding image URIs:", err)
		return
	}

	// Convert to JSON and return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(card)
}

// Get card by name
func getCardByName(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Extract card name from URL
	cardName := r.URL.Path[len("/card/"):]
	if cardName == "" {
		http.Error(w, "Card name is required", http.StatusBadRequest)
		return
	}

	query := `
		SELECT name, mana_cost, image_uris, type_line, oracle_text, set, set_name, set_uri, set_id, set_type, set_search_uri, scryfall_set_uri
		FROM oracle_cards
		WHERE LOWER(name) = LOWER($1)
		LIMIT 1;
	`
	var card handlers.OracleCard
	var imageURIsJSON []byte

	err := db.QueryRow(query, cardName).Scan(
		&card.Name, &card.ManaCost, &imageURIsJSON, &card.TypeLine, &card.OracleText,
		&card.Set, &card.SetName, &card.SetURI, &card.SetID, &card.SetType,
		&card.SetSearchURI, &card.ScryfallSetURI,
	)
	if err == sql.ErrNoRows {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error fetching card", http.StatusInternalServerError)
		log.Println("❌ Error fetching card:", err)
		return
	}

	// Convert JSONB data from database into Go struct
	if err := json.Unmarshal(imageURIsJSON, &card.ImageURIs); err != nil {
		http.Error(w, "Error decoding image URIs", http.StatusInternalServerError)
		log.Println("❌ Error decoding image URIs:", err)
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
	connectDB()

	// Ensure tables exist
	ensureTables(db)

	// Start the bulk data fetch scheduler
	utils.StartScheduler(db)

	// API endpoints
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Mana Tomb API is running!"))
	})
	http.HandleFunc("/random-card", randomCard)
	http.HandleFunc("/card/", getCardByName) // New card search endpoint

	log.Println("🚀 Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
