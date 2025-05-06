package cards

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/quehorrifico/mana-tomb/backend/middleware"
	"github.com/quehorrifico/mana-tomb/backend/models"
)

var DB *sql.DB

// Get a random card from the database
func GetRandomCard(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCORS(w) // Add CORS headers
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Use the existing db connection
	query := `
		SELECT name, mana_cost, image_uris, type_line, oracle_text, set, set_name, set_uri, set_id, set_type, set_search_uri, scryfall_set_uri
		FROM oracle_cards
		ORDER BY random(), (SELECT COUNT(*) FROM jsonb_object_keys(image_uris)) DESC
		LIMIT 1;
	`

	var card models.OracleCard
	var imageURIsJSON []byte

	// Fetch a random card
	err := DB.QueryRow(query).Scan(&card.Name, &card.ManaCost, &imageURIsJSON, &card.TypeLine, &card.OracleText, &card.Set, &card.SetName, &card.SetURI, &card.SetID, &card.SetType, &card.SetSearchURI, &card.ScryfallSetURI)
	if err != nil {
		http.Error(w, "Error fetching card", http.StatusInternalServerError)
		log.Println("❌ Error fetching random card:", err)
		return
	}

	// Random card found, return it
	if err := json.Unmarshal(imageURIsJSON, &card.ImageURIs); err != nil {
		http.Error(w, "Error decoding image URIs", http.StatusInternalServerError)
		log.Println("❌ Error decoding image URIs:", err)
		return
	}

	// Send the single exact match with a flag
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"exact_match": card,
	})
}

// Get cards by fuzzy name search
func GetCardByName(w http.ResponseWriter, r *http.Request) {
	middleware.EnableCORS(w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	cardName := r.URL.Path[len("/card/"):]
	if cardName == "" {
		http.Error(w, "Card name is required", http.StatusBadRequest)
		return
	}

	// Try finding an **exact match** first
	queryExact := `
		SELECT name, mana_cost, image_uris, type_line, oracle_text, set, set_name, set_uri, set_id, set_type, set_search_uri, scryfall_set_uri
		FROM oracle_cards
		WHERE name ILIKE $1
		AND oracle_text IS NOT NULL AND oracle_text <> ''
		ORDER BY (SELECT COUNT(*) FROM jsonb_object_keys(image_uris)) DESC
		LIMIT 1;
	`
	var exactMatch models.OracleCard
	var imageURIsJSON []byte
	err := DB.QueryRow(queryExact, cardName).Scan(
		&exactMatch.Name, &exactMatch.ManaCost, &imageURIsJSON, &exactMatch.TypeLine, &exactMatch.OracleText,
		&exactMatch.Set, &exactMatch.SetName, &exactMatch.SetURI, &exactMatch.SetID, &exactMatch.SetType,
		&exactMatch.SetSearchURI, &exactMatch.ScryfallSetURI,
	)
	if err == nil {
		// Exact match found, return it
		if err := json.Unmarshal(imageURIsJSON, &exactMatch.ImageURIs); err != nil {
			http.Error(w, "Error decoding image URIs", http.StatusInternalServerError)
			log.Println("❌ Error decoding image URIs:", err)
			return
		}

		// Send the single exact match with a flag
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"exact_match": exactMatch,
		})
		return
	}

	// If no exact match, return **fuzzy matches**
	queryFuzzy := `
		SELECT name, mana_cost, image_uris, type_line, oracle_text, set, set_name, set_uri, set_id, set_type, set_search_uri, scryfall_set_uri
		FROM oracle_cards
		WHERE name ILIKE '%' || $1 || '%'
		LIMIT 10;
	`
	rows, err := DB.Query(queryFuzzy, cardName)
	if err != nil {
		http.Error(w, "Error fetching cards", http.StatusInternalServerError)
		log.Println("❌ Error fetching cards:", err)
		return
	}
	defer rows.Close()

	var cards []models.OracleCard
	for rows.Next() {
		var card models.OracleCard
		err := rows.Scan(
			&card.Name, &card.ManaCost, &imageURIsJSON, &card.TypeLine, &card.OracleText,
			&card.Set, &card.SetName, &card.SetURI, &card.SetID, &card.SetType,
			&card.SetSearchURI, &card.ScryfallSetURI,
		)
		if err != nil {
			log.Println("❌ Error scanning card:", err)
			continue
		}

		// Convert JSONB data from database into Go struct
		if err := json.Unmarshal(imageURIsJSON, &card.ImageURIs); err != nil {
			log.Println("❌ Error decoding image URIs:", err)
			continue
		}

		cards = append(cards, card)
	}

	// If no fuzzy matches found
	if len(cards) == 0 {
		http.Error(w, "No cards found", http.StatusNotFound)
		return
	}

	// Convert to JSON and return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"fuzzy_matches": cards,
	})
}
