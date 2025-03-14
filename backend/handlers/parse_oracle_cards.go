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

type OracleCardImageURIs struct {
	Small      string `json:"small"`
	Normal     string `json:"normal"`
	Large      string `json:"large"`
	PNG        string `json:"png"`
	ArtCrop    string `json:"art_crop"`
	BorderCrop string `json:"border_crop"`
}

type OracleCardLegalities struct {
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

type OracleCardPrices struct {
	USD       string `json:"usd"`
	USDFoil   string `json:"usd_foil"`
	USDEtched string `json:"usd_etched"`
	EUR       string `json:"eur"`
	EURFoil   string `json:"eur_foil"`
	Tix       string `json:"tix"`
}

type OracleCard struct {
	ID              string               `json:"id"`
	OracleID        string               `json:"oracle_id"`
	MultiverseIDs   []int                `json:"multiverse_ids"`
	MTGOID          int                  `json:"mtgo_id"`
	MTGOFoilID      int                  `json:"mtgo_foil_id"`
	TCGPlayerID     int                  `json:"tcgplayer_id"`
	CardMarketID    int                  `json:"cardmarket_id"`
	Name            string               `json:"name"`
	Lang            string               `json:"lang"`
	ReleasedAt      string               `json:"released_at"`
	URI             string               `json:"uri"`
	ScryfallURI     string               `json:"scryfall_uri"`
	Layout          string               `json:"layout"`
	HighResImage    bool                 `json:"highres_image"`
	ImageStatus     string               `json:"image_status"`
	ImageURIs       OracleCardImageURIs  `json:"image_uris"`
	ManaCost        string               `json:"mana_cost"`
	CMC             float64              `json:"cmc"`
	TypeLine        string               `json:"type_line"`
	OracleText      string               `json:"oracle_text"`
	Colors          []string             `json:"colors"`
	ColorIdentity   []string             `json:"color_identity"`
	Keywords        []string             `json:"keywords"`
	Legalities      OracleCardLegalities `json:"legalities"`
	Games           []string             `json:"games"`
	Reserved        bool                 `json:"reserved"`
	GameChanger     bool                 `json:"game_changer"`
	Foil            bool                 `json:"foil"`
	NonFoil         bool                 `json:"nonfoil"`
	Finishes        []string             `json:"finishes"`
	Oversized       bool                 `json:"oversized"`
	Promo           bool                 `json:"promo"`
	Reprint         bool                 `json:"reprint"`
	Variation       bool                 `json:"variation"`
	SetID           string               `json:"set_id"`
	Set             string               `json:"set"`
	SetName         string               `json:"set_name"`
	SetType         string               `json:"set_type"`
	SetURI          string               `json:"set_uri"`
	SetSearchURI    string               `json:"set_search_uri"`
	ScryfallSetURI  string               `json:"scryfall_set_uri"`
	RulingsURI      string               `json:"rulings_uri"`
	PrintsSearchURI string               `json:"prints_search_uri"`
	CollectorNumber string               `json:"collector_number"`
	Digital         bool                 `json:"digital"`
	Rarity          string               `json:"rarity"`
	FlavorText      string               `json:"flavor_text"`
	CardBackID      string               `json:"card_back_id"`
	Artist          string               `json:"artist"`
	ArtistIDs       []string             `json:"artist_ids"`
	IllustrationID  string               `json:"illustration_id"`
	BorderColor     string               `json:"border_color"`
	Frame           string               `json:"frame"`
	FullArt         bool                 `json:"full_art"`
	Textless        bool                 `json:"textless"`
	Booster         bool                 `json:"booster"`
	StorySpotlight  bool                 `json:"story_spotlight"`
	EDHRecRank      int                  `json:"edhrec_rank"`
	Prices          OracleCardPrices     `json:"prices"`
}

// ParseAndStoreOracleCards fetches JSON from the download_uri and stores it in new tables
func ParseAndStoreOracleCards(db *sql.DB) {
	log.Println("🔄 Fetching JSON from stored bulk data download URIs...")

	// Query database for all bulk data items
	// rows, err := db.Query("SELECT type, name, download_uri FROM bulk_data")
	// if err != nil {
	// 	log.Printf("❌ Error querying bulk_data table: %v\n", err)
	// 	return
	// }
	// defer rows.Close()

	// Query database for all bulk data items with type "oracle_cards"
	rows, err := db.Query("SELECT type, name, download_uri FROM bulk_data WHERE type = 'oracle_cards'")
	if err != nil {
		log.Printf("❌ Error querying bulk_data table: %v\n", err)
		return
	}
	defer rows.Close()

	// Iterate over oracle_cards bulk data item
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

		// Create oracle_cards table if it does not exist
		err = createOracleCardsTableForBulkData(db, dataType)
		if err != nil {
			log.Printf("❌ Error creating table for %s: %v\n", tableName, err)
			continue
		}

		// Parse JSON
		var oracle_cards []OracleCard
		if err := json.Unmarshal(body, &oracle_cards); err != nil {
			log.Printf("❌ Error unmarshaling JSON for %s: %v\n", tableName, err)
			continue
		}

		// Insert parsed data into the appropriate table
		err = insertCardsIntoOracleCardsTable(db, tableName, oracle_cards)
		if err != nil {
			log.Printf("❌ Error inserting cards into %s: %v\n", tableName, err)
		} else {
			log.Printf("✅ Successfully inserted %d cards into %s\n", len(oracle_cards), tableName)
		}
	}
}

func createOracleCardsTableForBulkData(db *sql.DB, tableName string) error {
	safeTableName := sanitizeOracleCardsTableName(tableName)

	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id UUID PRIMARY KEY,
		oracle_id UUID,
		multiverse_ids INTEGER[],
		mtgo_id INTEGER,
		mtgo_foil_id INTEGER,
		tcgplayer_id INTEGER,
		cardmarket_id INTEGER,
		name TEXT NOT NULL,
		lang TEXT NOT NULL,
		released_at DATE,
		uri TEXT NOT NULL,
		scryfall_uri TEXT NOT NULL,
		layout TEXT NOT NULL,
		highres_image BOOLEAN,
		image_status TEXT,
		image_uris JSONB,
		mana_cost TEXT,
		cmc NUMERIC,
		type_line TEXT NOT NULL,
		oracle_text TEXT,
		power TEXT,
		toughness TEXT,
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
		UNIQUE(id)
	);`, safeTableName)

	_, err := db.Exec(query)
	return err
}

// insertCardsIntoTable inserts parsed JSON data into the correct table
func insertCardsIntoOracleCardsTable(db *sql.DB, tableName string, cards []OracleCard) error {
	safeTableName := sanitizeOracleCardsTableName(tableName)
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
	`, safeTableName))
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

// sanitizeTableName makes sure the table name is safe for SQL
func sanitizeOracleCardsTableName(name string) string {
	return strings.ReplaceAll(strings.ToLower(name), " ", "_")
}
