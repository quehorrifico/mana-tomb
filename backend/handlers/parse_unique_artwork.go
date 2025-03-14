package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/lib/pq"
)

type UniqueArtworkImageURIs struct {
	Small      string `json:"small"`
	Normal     string `json:"normal"`
	Large      string `json:"large"`
	PNG        string `json:"png"`
	ArtCrop    string `json:"art_crop"`
	BorderCrop string `json:"border_crop"`
}

type UniqueArtworkLegalities struct {
	Standard        string `json:"standard"`
	Future          string `json:"future"`
	Historic        string `json:"historic"`
	Timeless        string `json:"timeless"`
	Gladiator       string `json:"gladiator"`
	Pioneer         string `json:"pioneer"`
	Explorer        string `json:"explorer"`
	Modern          string `json:"modern"`
	Legacy          string `json:"legacy"`
	Pauper          string `json:"pauper"`
	Vintage         string `json:"vintage"`
	Penny           string `json:"penny"`
	Commander       string `json:"commander"`
	Oathbreaker     string `json:"oathbreaker"`
	StandardBrawl   string `json:"standardbrawl"`
	Brawl           string `json:"brawl"`
	Alchemy         string `json:"alchemy"`
	PauperCommander string `json:"paupercommander"`
	Duel            string `json:"duel"`
	Oldschool       string `json:"oldschool"`
	Premodern       string `json:"premodern"`
	Predh           string `json:"predh"`
}

type UniqueArtworkPrices struct {
	USD       string `json:"usd"`
	USDFoil   string `json:"usd_foil"`
	USDEtched string `json:"usd_etched"`
	EUR       string `json:"eur"`
	EURFoil   string `json:"eur_foil"`
	Tix       string `json:"tix"`
}

type UniqueArtworkRelatedURIs struct {
	Gatherer                  string `json:"gatherer"`
	TCGPlayerInfiniteArticles string `json:"tcgplayer_infinite_articles"`
	TCGPlayerInfiniteDecks    string `json:"tcgplayer_infinite_decks"`
	EDHRec                    string `json:"edhrec"`
}

type UniqueArtworkPurchaseURIs struct {
	TCGPlayer   string `json:"tcgplayer"`
	Cardmarket  string `json:"cardmarket"`
	CardHoarder string `json:"cardhoarder"`
}

type UniqueArtworkCards struct {
	Object          string                   `json:"object"`
	ID              string                   `json:"id"`
	OracleID        string                   `json:"oracle_id"`
	MultiverseIDs   []int                    `json:"multiverse_ids"`
	MTGOID          int                      `json:"mtgo_id"`
	MTGOFoilID      int                      `json:"mtgo_foil_id"`
	TCGPlayerID     int                      `json:"tcgplayer_id"`
	CardMarketID    int                      `json:"cardmarket_id"`
	ArenaID         int                      `json:"arena_id"`
	Name            string                   `json:"name"`
	Lang            string                   `json:"lang"`
	ReleasedAt      string                   `json:"released_at"`
	URI             string                   `json:"uri"`
	ScryfallURI     string                   `json:"scryfall_uri"`
	Layout          string                   `json:"layout"`
	HighResImage    bool                     `json:"highres_image"`
	ImageStatus     string                   `json:"image_status"`
	ImageURIs       UniqueArtworkImageURIs   `json:"image_uris"`
	ManaCost        string                   `json:"mana_cost"`
	CMC             float64                  `json:"cmc"`
	TypeLine        string                   `json:"type_line"`
	OracleText      string                   `json:"oracle_text"`
	Colors          []string                 `json:"colors"`
	ColorIdentity   []string                 `json:"color_identity"`
	Keywords        []string                 `json:"keywords"`
	Legalities      UniqueArtworkLegalities  `json:"legalities"`
	Games           []string                 `json:"games"`
	Reserved        bool                     `json:"reserved"`
	GameChanger     bool                     `json:"game_changer"`
	Foil            bool                     `json:"foil"`
	NonFoil         bool                     `json:"nonfoil"`
	Finishes        []string                 `json:"finishes"`
	Oversized       bool                     `json:"oversized"`
	Promo           bool                     `json:"promo"`
	Reprint         bool                     `json:"reprint"`
	Variation       bool                     `json:"variation"`
	SetID           string                   `json:"set_id"`
	Set             string                   `json:"set"`
	SetName         string                   `json:"set_name"`
	SetType         string                   `json:"set_type"`
	SetURI          string                   `json:"set_uri"`
	SetSearchURI    string                   `json:"set_search_uri"`
	ScryfallSetURI  string                   `json:"scryfall_set_uri"`
	RulingsURI      string                   `json:"rulings_uri"`
	PrintsSearchURI string                   `json:"prints_search_uri"`
	CollectorNumber string                   `json:"collector_number"`
	Digital         bool                     `json:"digital"`
	Rarity          string                   `json:"rarity"`
	FlavorText      string                   `json:"flavor_text"`
	CardBackID      string                   `json:"card_back_id"`
	Artist          string                   `json:"artist"`
	ArtistIDs       []string                 `json:"artist_ids"`
	IllustrationID  string                   `json:"illustration_id"`
	BorderColor     string                   `json:"border_color"`
	Frame           string                   `json:"frame"`
	FullArt         bool                     `json:"full_art"`
	Textless        bool                     `json:"textless"`
	Booster         bool                     `json:"booster"`
	StorySpotlight  bool                     `json:"story_spotlight"`
	EDHRecRank      int                      `json:"edhrec_rank"`
	PennyRank       int                      `json:"penny_rank"`
	Prices          UniqueArtworkPrices      `json:"prices"`
	RelatedURIs     UniqueArtworkRelatedURIs `json:"related_uris"`
}

// ParseAndStoreUniqueArtwork fetches JSON from the download_uri and stores it in new tables
func ParseAndStoreUniqueArtwork(db *sql.DB) {
	log.Println("🔄 Fetching JSON from stored bulk data download URIs...")

	// Query database for all bulk data items with type "unique_artwork"
	rows, err := db.Query("SELECT type, name, download_uri FROM bulk_data WHERE type = 'unique_artwork'")
	if err != nil {
		log.Printf("❌ Error querying bulk_data table: %v\n", err)
		return
	}
	defer rows.Close()

	// Iterate over unique_artwork bulk data item
	for rows.Next() {
		var dataType, tableName, downloadURI string
		if err := rows.Scan(&dataType, &tableName, &downloadURI); err != nil {
			log.Printf("❌ Error scanning row: %v\n", err)
			continue
		}

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

		// Create unique_artwork table if it does not exist
		err = createUniqueArtworkTableForBulkData(db, dataType)
		if err != nil {
			log.Printf("❌ Error creating table for %s: %v\n", tableName, err)
			continue
		}

		// Parse JSON
		var unique_artwork_cards []UniqueArtworkCards
		if err := json.Unmarshal(body, &unique_artwork_cards); err != nil {
			log.Printf("❌ Error unmarshaling JSON for %s: %v\n", tableName, err)
			continue
		}

		// Insert parsed data into the appropriate table
		err = insertCardsIntoUniqueArtworkTable(db, tableName, unique_artwork_cards)
		if err != nil {
			log.Printf("❌ Error inserting cards into %s: %v\n", tableName, err)
		} else {
			log.Printf("✅ Successfully inserted %d cards into %s\n", len(unique_artwork_cards), tableName)
		}
	}
}

// createUniqueArtworkTableForBulkData creates a unique_artwork table with the necessary fields
func createUniqueArtworkTableForBulkData(db *sql.DB, tableName string) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id TEXT PRIMARY KEY,
			object TEXT,
			oracle_id TEXT,
			multiverse_ids INTEGER[],
			mtgo_id INTEGER,
			mtgo_foil_id INTEGER,
			tcgplayer_id INTEGER,
			cardmarket_id INTEGER,
			arena_id INTEGER,
			name TEXT,
			lang TEXT,
			released_at DATE,
			uri TEXT,
			scryfall_uri TEXT,
			layout TEXT,
			highres_image BOOLEAN,
			image_status TEXT,
			image_uris JSONB,
			mana_cost TEXT,
			cmc FLOAT,
			type_line TEXT,
			oracle_text TEXT,
			colors TEXT[],
			color_identity TEXT[],
			keywords TEXT[],
			legalities JSONB,
			games TEXT[],
			reserved BOOLEAN,
			game_changer BOOLEAN,
			foil BOOLEAN,
			nonfoil BOOLEAN,
			finishes TEXT[],
			oversized BOOLEAN,
			promo BOOLEAN,
			reprint BOOLEAN,
			variation BOOLEAN,
			set_id TEXT,
			set TEXT,
			set_name TEXT,
			set_type TEXT,
			set_uri TEXT,
			set_search_uri TEXT,
			scryfall_set_uri TEXT,
			rulings_uri TEXT,
			prints_search_uri TEXT,
			collector_number TEXT,
			digital BOOLEAN,
			rarity TEXT,
			flavor_text TEXT,
			card_back_id TEXT,
			artist TEXT,
			artist_ids TEXT[],
			illustration_id TEXT,
			border_color TEXT,
			frame TEXT,
			full_art BOOLEAN,
			textless BOOLEAN,
			booster BOOLEAN,
			story_spotlight BOOLEAN,
			edhrec_rank INTEGER,
			penny_rank INTEGER,
			prices JSONB,
			related_uris JSONB
		);
	`, tableName)

	_, err := db.Exec(query)
	return err
}

// insertCardsIntoUniqueArtworkTable inserts parsed JSON data into the correct table
func insertCardsIntoUniqueArtworkTable(db *sql.DB, tableName string, cards []UniqueArtworkCards) error {
	safeTableName := sanitizeUniqueArtworkTableName(tableName)
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
	`, safeTableName))
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
func sanitizeUniqueArtworkTableName(name string) string {
	return strings.ReplaceAll(strings.ToLower(name), " ", "_")
}
