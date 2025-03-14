package utils

import (
	"database/sql"
	"log"
	"time"
	// "github.com/quehorrifico/mana-tomb/handlers"
)

// StartScheduler runs FetchAndStoreBulkData once per day
func StartScheduler(db *sql.DB) {
	go func() {
		for {
			log.Println("🔄 Fetching and storing bulk data from Scryfall...")
			// handlers.FetchAndStoreBulkData(db)
			// handlers.ParseAndStoreOracleCards(db)
			// handlers.ParseAndStoreUniqueArtwork(db)

			log.Println("⏳ Next bulk data update in 24 hours...")
			time.Sleep(24 * time.Hour) // Runs once every 24 hours
		}
	}()
}
