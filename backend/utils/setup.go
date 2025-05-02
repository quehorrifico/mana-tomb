package utils

import (
	"database/sql"
	"fmt"
	"log"
)

// EnsureUserTable creates the users table if it doesn't exist.
func EnsureUserTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}
	log.Println("✅ users table checked/created successfully!")
	return nil
}

// EnsureBulkDataTable creates the bulk_data table if it doesn't exist.
func EnsureBulkDataTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS bulk_data (
		id UUID PRIMARY KEY,
		type TEXT,
		updated_at TIMESTAMPTZ,
		download_uri TEXT,
		name TEXT,
		description TEXT,
		size BIGINT,
		content_type TEXT,
		content_encoding TEXT
	);
	`
	_, err := db.Exec(query)
	return err
}

// EnsureOracleCardsTable creates the oracle_cards table if it doesn't exist.
func EnsureOracleCardsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS oracle_cards (
		id UUID PRIMARY KEY,
		name TEXT,
		mana_cost TEXT,
		type_line TEXT,
		oracle_text TEXT,
		image_uris JSONB,
		set TEXT,
		set_name TEXT
	);`
	_, err := db.Exec(query)
	return err
}

// EnsureUniqueArtworkTable creates the unique_artwork table if it doesn't exist.
func EnsureUniqueArtworkTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS unique_artwork (
		id UUID PRIMARY KEY,
		name TEXT,
		image_uris JSONB,
		set TEXT,
		set_name TEXT
	);`
	_, err := db.Exec(query)
	return err
}
