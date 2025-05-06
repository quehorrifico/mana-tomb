package decks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lib/pq"
	"github.com/quehorrifico/mana-tomb/backend/db"
	"github.com/quehorrifico/mana-tomb/backend/models"
)

var DB *sql.DB

func GetDecksByUser(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value("userID")
	userID, ok := userIDRaw.(int)
	if !ok || userID == 0 {
		http.Error(w, "Missing or invalid user ID", http.StatusUnauthorized)
		return
	}

	rows, err := DB.Query(`SELECT id, user_id, name, description, commander, cards FROM proto_commander_decks WHERE user_id = $1`, userID)
	if err != nil {
		http.Error(w, "Failed to fetch decks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var commander_decks []models.ProtoCommanderDeck
	for rows.Next() {
		var commander_deck models.ProtoCommanderDeck
		if err := rows.Scan(&commander_deck.DeckID, &commander_deck.UserID, &commander_deck.Name, &commander_deck.Description, &commander_deck.Commander, pq.Array(&commander_deck.Cards)); err != nil {
			fmt.Printf("Error scanning commander deck: %v\n", err)
			http.Error(w, "Error scanning commander deck", http.StatusInternalServerError)
			return
		}
		commander_decks = append(commander_decks, commander_deck)
	}

	json.NewEncoder(w).Encode(commander_decks)
}

func CreateDeck(w http.ResponseWriter, r *http.Request) {
	userIDRaw := r.Context().Value("userID")
	userID, ok := userIDRaw.(int)
	if !ok || userID == 0 {
		http.Error(w, "Missing or invalid user ID", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var commander_deck models.ProtoCommanderDeck
	if err := json.NewDecoder(r.Body).Decode(&commander_deck); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var deckID int
	err = tx.QueryRow(`
		INSERT INTO proto_commander_decks (user_id, name, description, commander, cards)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`, userID, commander_deck.Name, commander_deck.Description, commander_deck.Commander, pq.Array(commander_deck.Cards)).Scan(&deckID)
	if err != nil {
		http.Error(w, "Failed to create deck", http.StatusInternalServerError)
		return
	}
	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"deck_id": deckID,
	})
}

func GetDeckByID(w http.ResponseWriter, r *http.Request) {
	// parts := strings.Split(r.URL.Path, "/")
	// if len(parts) < 3 {
	// 	http.Error(w, "Invalid deck ID in path", http.StatusBadRequest)
	// 	return
	// }
	// id := parts[2]
	id := r.URL.Path[len("/decks/"):]
	if id == "" {
		http.Error(w, "Invalid deck ID in path", http.StatusBadRequest)
		return
	}

	var commander_deck models.ProtoCommanderDeck
	query := `SELECT * FROM proto_commander_decks WHERE id = $1`
	err := DB.QueryRow(query, id).Scan(&commander_deck.DeckID, &commander_deck.UserID, &commander_deck.Name, &commander_deck.Description, &commander_deck.Commander, pq.Array(&commander_deck.Cards))
	if err == sql.ErrNoRows {
		http.Error(w, "Commander deck not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to retrieve commander deck", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Deck ID: %d, User ID: %d, Name: %s, Description: %s, Commander: %s, Cards: %v\n", commander_deck.DeckID, commander_deck.UserID, commander_deck.Name, commander_deck.Description, commander_deck.Commander, commander_deck.Cards)
	json.NewEncoder(w).Encode(commander_deck)
}
