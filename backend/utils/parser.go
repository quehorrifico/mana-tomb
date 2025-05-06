package utils

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"encoding/json"
	"io"
	"net/http"

	"github.com/lib/pq"
	"github.com/quehorrifico/mana-tomb/backend/models"
)

// FetchAndParseBulkData fetches bulk data from Scryfall and stores it in PostgreSQL (oracle_cards, unique_artwork, etc.)
func FetchAndParseBulkData(db *sql.DB) {
	log.Println("🔄 Fetching and storing bulk data items from Scryfall...")
	resp, err := http.Get("https://api.scryfall.com/bulk-data")
	if err != nil {
		log.Printf("Error fetching bulk data: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
		return
	}

	var bulkDataResponse models.BulkDataResponse

	err = json.Unmarshal(body, &bulkDataResponse)
	if err != nil {
		log.Printf("Error unmarshaling JSON: %v\n", err)
		return
	}

	// Now we have the correct slice of bulk data items
	bulkDataList := bulkDataResponse.Data

	// Insert/update data in PostgreSQL
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v\n", err)
		return
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO bulk_data (id, uri, type, name, description, download_uri, updated_at, size, content_type, content_encoding)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO UPDATE
		SET uri = EXCLUDED.uri,
		    type = EXCLUDED.type,
		    name = EXCLUDED.name,
		    description = EXCLUDED.description,
		    download_uri = EXCLUDED.download_uri,
		    updated_at = EXCLUDED.updated_at,
		    size = EXCLUDED.size,
		    content_type = EXCLUDED.content_type,
		    content_encoding = EXCLUDED.content_encoding;
	`)
	if err != nil {
		log.Printf("Error preparing statement: %v\n", err)
		return
	}
	defer stmt.Close()

	for _, data := range bulkDataList {
		_, err := stmt.Exec(data.ID, data.URI, data.Type, data.Name, data.Description, data.DownloadURI, data.UpdatedAt, data.Size, data.ContentType, data.ContentEncoding)
		if err != nil {
			log.Printf("Error inserting bulk data: %v\n", err)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v\n", err)
		return
	}

	log.Println("Successfully updated bulk data.")
}

// FetchAndParseOracleCards fetches JSON from the download_uri and stores it in new tables
func FetchAndParseOracleCards(db *sql.DB) {
	log.Println("📥 Parsing oracle cards from bulk data...")

	// Query database for all bulk data items with type "oracle_cards"
	rows, err := db.Query("SELECT type, name, download_uri FROM bulk_data WHERE type = 'oracle_cards'")
	if err != nil {
		log.Printf("❌ Error querying bulk_data table: %v\n", err)
		return
	}
	defer rows.Close()

	// Iterate over oracle_cards bulk data item
	for rows.Next() {
		var dataType, unsafeTableName, downloadURI string
		if err := rows.Scan(&dataType, &unsafeTableName, &downloadURI); err != nil {
			log.Printf("❌ Error scanning row: %v\n", err)
			continue
		}
		tableName := sanitizeTableName(unsafeTableName)

		log.Printf("📥 Fetching JSON for %s from %s\n", tableName, downloadURI)

		// Fetch the JSON payload
		resp, err := http.Get(downloadURI)
		if err != nil {
			log.Printf("❌ Failed to fetch JSON from %s: %v\n", downloadURI, err)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("❌ Error reading JSON response: %v\n", err)
			continue
		}

		// Create oracle_cards table if it does not exist
		err = EnsureOracleCardsTable(db)
		if err != nil {
			log.Printf("❌ Error creating table for %s: %v\n", tableName, err)
			continue
		}

		// Parse JSON
		var oracleCards []models.OracleCard
		if err := json.Unmarshal(body, &oracleCards); err != nil {
			log.Printf("❌ Error unmarshaling JSON for %s: %v\n", tableName, err)
			continue
		}

		// Insert parsed data into the appropriate table
		err = InsertOracleCards(db, tableName, oracleCards)
		if err != nil {
			log.Printf("❌ Error inserting cards into %s: %v\n", tableName, err)
		} else {
			log.Printf("✅ Successfully inserted %d cards into %s\n", len(oracleCards), tableName)
		}
	}
}

// ParseAndParseUniqueArtwork fetches JSON from the download_uri and stores it in new tables
func ParseAndParseUniqueArtwork(db *sql.DB) {
	log.Println("📥 Parsing unique artwork from bulk data...")

	// Query database for all bulk data items with type "oracle_cards"
	rows, err := db.Query("SELECT type, name, download_uri FROM bulk_data WHERE type = 'oracle_cards'")
	if err != nil {
		log.Printf("❌ Error querying bulk_data table: %v\n", err)
		return
	}
	defer rows.Close()

	// Iterate over oracle_cards bulk data item
	for rows.Next() {
		var dataType, unsafeTableName, downloadURI string
		if err := rows.Scan(&dataType, &unsafeTableName, &downloadURI); err != nil {
			log.Printf("❌ Error scanning row: %v\n", err)
			continue
		}
		tableName := sanitizeTableName(unsafeTableName)

		log.Printf("📥 Fetching JSON for %s from %s\n", tableName, downloadURI)

		// Fetch the JSON payload
		resp, err := http.Get(downloadURI)
		if err != nil {
			log.Printf("❌ Failed to fetch JSON from %s: %v\n", downloadURI, err)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("❌ Error reading JSON response: %v\n", err)
			continue
		}

		// Create oracle_cards table if it does not exist
		err = EnsureOracleCardsTable(db)
		if err != nil {
			log.Printf("❌ Error creating table for %s: %v\n", tableName, err)
			continue
		}

		// Parse JSON
		var oracleCards []models.UniqueArtworkCard
		if err := json.Unmarshal(body, &oracleCards); err != nil {
			log.Printf("❌ Error unmarshaling JSON for %s: %v\n", tableName, err)
			continue
		}

		// Insert parsed data into the appropriate table
		err = InsertUniqueArtwork(db, tableName, oracleCards)
		if err != nil {
			log.Printf("❌ Error inserting cards into %s: %v\n", tableName, err)
		} else {
			log.Printf("✅ Successfully inserted %d cards into %s\n", len(oracleCards), tableName)
		}
	}
}

// InsertOracleCards inserts parsed JSON data into the correct table
func InsertOracleCards(db *sql.DB, tableName string, cards []models.OracleCard) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(fmt.Sprintf(`
		INSERT INTO %s (
			id, oracle_id, multiverse_ids, mtgo_id, mtgo_foil_id, tcgplayer_id, cardmarket_id,
			name, lang, released_at, uri, scryfall_uri, layout, highres_image, image_status, image_uris,
			mana_cost, cmc, type_line, oracle_text, colors, color_identity, keywords, legalities,
			games, reserved, game_changer, foil, nonfoil, finishes, oversized, promo, reprint, variation,
			set_id, set, set_name, set_type, set_uri, set_search_uri, scryfall_set_uri, rulings_uri,
			prints_search_uri, collector_number, digital, rarity, flavor_text, card_back_id, artist,
			artist_ids, illustration_id, border_color, frame, full_art, textless, booster,
			story_spotlight, edhrec_rank, prices
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21, $22, $23, $24,
			$25, $26, $27, $28, $29, $30, $31, $32, $33, $34,
			$35, $36, $37, $38, $39, $40, $41, $42,
			$43, $44, $45, $46, $47, $48, $49,
			$50, $51, $52, $53, $54, $55, $56,
			$57, $58, $59
		)
		ON CONFLICT (id) DO NOTHING;
	`, tableName))
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	for _, card := range cards {
		imageURIsJSON, _ := json.Marshal(card.ImageURIs)
		legalitiesJSON, _ := json.Marshal(card.Legalities)
		pricesJSON, _ := json.Marshal(card.Prices)

		args := []interface{}{
			card.ID, card.OracleID, pq.Array(card.MultiverseIDs), card.MTGOID, card.MTGOFoilID, card.TCGPlayerID, card.CardMarketID,
			card.Name, card.Lang, card.ReleasedAt, card.URI, card.ScryfallURI, card.Layout, card.HighResImage, card.ImageStatus, imageURIsJSON,
			card.ManaCost, card.CMC, card.TypeLine, card.OracleText, pq.Array(card.Colors), pq.Array(card.ColorIdentity), pq.Array(card.Keywords), legalitiesJSON,
			pq.Array(card.Games), card.Reserved, card.GameChanger, card.Foil, card.NonFoil, pq.Array(card.Finishes), card.Oversized, card.Promo, card.Reprint, card.Variation,
			card.SetID, card.Set, card.SetName, card.SetType, card.SetURI, card.SetSearchURI, card.ScryfallSetURI, card.RulingsURI,
			card.PrintsSearchURI, card.CollectorNumber, card.Digital, card.Rarity, card.FlavorText, card.CardBackID, card.Artist,
			pq.Array(card.ArtistIDs), card.IllustrationID, card.BorderColor, card.Frame, card.FullArt, card.Textless, card.Booster,
			card.StorySpotlight, card.EDHRecRank, pricesJSON,
		}
		_, err = stmt.Exec(args...)
		if err != nil {
			log.Printf("⚠️ Error inserting card %s: %v", card.Name, err)
			continue
		}
	}

	err = tx.Commit()
	return err
}

// InsertUniqueArtwork inserts parsed JSON data into the correct table
func InsertUniqueArtwork(db *sql.DB, tableName string, cards []models.UniqueArtworkCard) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(fmt.Sprintf(`
		INSERT INTO %s (
			id, object, oracle_id, multiverse_ids, mtgo_id, mtgo_foil_id, tcgplayer_id, cardmarket_id, arena_id,
			name, lang, released_at, uri, scryfall_uri, layout, highres_image, image_status, image_uris,
			mana_cost, cmc, type_line, oracle_text, colors, color_identity, keywords, legalities,
			games, reserved, game_changer, foil, nonfoil, finishes, oversized, promo, reprint, variation,
			set_id, set, set_name, set_type, set_uri, set_search_uri, scryfall_set_uri, rulings_uri,
			prints_search_uri, collector_number, digital, rarity, flavor_text, card_back_id, artist,
			artist_ids, illustration_id, border_color, frame, full_art, textless, booster,
			story_spotlight, edhrec_rank, penny_rank, prices, related_uris
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9,
			$10, $11, $12, $13, $14, $15, $16, $17, $18,
			$19, $20, $21, $22, $23, $24, $25, $26,
			$27, $28, $29, $30, $31, $32, $33, $34, $35,
			$36, $37, $38, $39, $40, $41, $42, $43,
			$44, $45, $46, $47, $48, $49, $50,
			$51, $52, $53, $54, $55, $56, $57,
			$58, $59, $60, $61, $62, $63
		)
		ON CONFLICT (id) DO NOTHING;
	`, tableName))
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	for _, card := range cards {
		// Convert JSON struct fields to JSONB
		imageURIsJSON, _ := json.Marshal(card.ImageURIs)
		legalitiesJSON, _ := json.Marshal(card.Legalities)
		pricesJSON, _ := json.Marshal(card.Prices)
		relatedURIsJSON, _ := json.Marshal(card.RelatedURIs)

		args := []interface{}{
			card.ID, card.Object, card.OracleID, pq.Array(card.MultiverseIDs), card.MTGOID, card.MTGOFoilID, card.TCGPlayerID, card.CardMarketID, card.ArenaID,
			card.Name, card.Lang, card.ReleasedAt, card.URI, card.ScryfallURI, card.Layout, card.HighResImage, card.ImageStatus, imageURIsJSON,
			card.ManaCost, card.CMC, card.TypeLine, card.OracleText, pq.Array(card.Colors), pq.Array(card.ColorIdentity), pq.Array(card.Keywords), legalitiesJSON,
			pq.Array(card.Games), card.Reserved, card.GameChanger, card.Foil, card.NonFoil, pq.Array(card.Finishes), card.Oversized, card.Promo, card.Reprint, card.Variation,
			card.SetID, card.Set, card.SetName, card.SetType, card.SetURI, card.SetSearchURI, card.ScryfallSetURI, card.RulingsURI,
			card.PrintsSearchURI, card.CollectorNumber, card.Digital, card.Rarity, card.FlavorText, card.CardBackID, card.Artist,
			pq.Array(card.ArtistIDs), card.IllustrationID, card.BorderColor, card.Frame, card.FullArt, card.Textless, card.Booster,
			card.StorySpotlight, card.EDHRecRank, card.PennyRank, pricesJSON, relatedURIsJSON,
		}
		_, err = stmt.Exec(args...)
		if err != nil {
			log.Printf("⚠️ Error inserting card %s: %v", card.Name, err)
			continue
		}
	}

	err = tx.Commit()
	return err
}

// sanitizeTableName makes sure the table name is safe for SQL
func sanitizeTableName(name string) string {
	return strings.ReplaceAll(strings.ToLower(name), " ", "_")
}
