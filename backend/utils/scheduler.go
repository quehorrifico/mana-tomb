package utils

import (
	"database/sql"
	"log"
	"time"

	"github.com/quehorrifico/mana-tomb/backend/db"
)

// StartScheduler runs FetchAndStoreBulkData once per day
func StartScheduler(db *sql.DB) {
	SetupDatabase()
	go func() {
		for {
			doFetch := false
			if !doFetch {
				log.Println("🛑 Fetching and storing process blocked, using existing data")
			} else {
				FetchAndParseBulkData(db)
				FetchAndParseOracleCards(db)
				ParseAndParseUniqueArtwork(db)
			}
			log.Println("⏳ Next bulk data update in 24 hours...")
			time.Sleep(24 * time.Hour) // Runs once every 24 hours
		}
	}()
}

// SetupDatabase ensures all domain tables exist.
func SetupDatabase() {
	db := db.GetDB()

	// Setup account tables
	if err := EnsureUserTable(db); err != nil {
		log.Fatalf("❌ Failed to create users table: %v", err)
	}
	if err := EnsureSessionsTable(db); err != nil {
		log.Fatalf("❌ Failed to create sessions table: %v", err)
	}

	// Setup card-related tables
	if err := EnsureBulkDataTable(db); err != nil {
		log.Fatalf("❌ Failed to create bulk_data table: %v", err)
	}
	if err := EnsureOracleCardsTable(db); err != nil {
		log.Fatalf("❌ Failed to create oracle_cards table: %v", err)
	}
	if err := EnsureUniqueArtworkTable(db); err != nil {
		log.Fatalf("❌ Failed to create unique_artwork table: %v", err)
	}
	if err := EnsureCommanderDecksTable(db); err != nil {
		log.Fatalf("❌ Failed to create proto_commander_decks table: %v", err)
	}

	log.Println("✅ All database tables initialized successfully")
}
