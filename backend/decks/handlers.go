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
		fmt.Println("here 3")
		http.Error(w, "Missing or invalid user ID", http.StatusUnauthorized)
		return
	}

	rows, err := DB.Query(`SELECT user_id, name, description, commander, cards FROM proto_commander_decks WHERE user_id = $1`, userID)
	if err != nil {
		http.Error(w, "Failed to fetch decks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var commander_decks []models.ProtoCommanderDeck
	for rows.Next() {
		var commander_deck models.ProtoCommanderDeck
		if err := rows.Scan(&commander_deck.UserID, &commander_deck.Name, &commander_deck.Description, &commander_deck.Commander, pq.Array(&commander_deck.Cards)); err != nil {
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
	fmt.Printf("userID: %+v\n", userID)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var commander_deck models.ProtoCommanderDeck
	fmt.Printf("request body: %+v\n", r.Body)
	if err := json.NewDecoder(r.Body).Decode(&commander_deck); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	fmt.Printf("commander_deck to insert: %+v\n", commander_deck)

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var deckID int
	fmt.Printf("userID: %+v\n", userID)
	fmt.Printf("commander_deck.Name: %+v\n", commander_deck.Name)
	fmt.Printf("commander_deck.Description: %+v\n", commander_deck.Description)
	fmt.Printf("commander_deck.Commander: %+v\n", commander_deck.Commander)
	fmt.Printf("commander_deck.Cards: %+v\n", commander_deck.Cards)
	err = tx.QueryRow(`
		INSERT INTO proto_commander_decks (user_id, name, description, commander, cards)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`, userID, commander_deck.Name, commander_deck.Description, commander_deck.Commander, pq.Array(commander_deck.Cards)).Scan(&deckID)
	if err != nil {
		fmt.Printf("Error inserting deck: %v\n", err)
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
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing deck ID", http.StatusBadRequest)
		return
	}

	var commander_deck models.ProtoCommanderDeck
	query := `SELECT user_id, name, description, commander, cards FROM proto_commander_decks WHERE id = $1`
	err := DB.QueryRow(query, id).Scan(&commander_deck.UserID, &commander_deck.Name, &commander_deck.Description, &commander_deck.Commander, &commander_deck.Cards)
	if err == sql.ErrNoRows {
		http.Error(w, "Commander deck not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to retrieve commander deck", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(commander_deck)
}
