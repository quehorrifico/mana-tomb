package utils

import (
	"database/sql"
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
	return err
}

// EnsureSessionsTable creates the sessions table if it doesn't exist.
func EnsureSessionsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS sessions (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		token TEXT UNIQUE NOT NULL,
		expires_at TIMESTAMPTZ NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(query)
	return err
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
		name TEXT NOT NULL,
		mana_cost TEXT,
		type_line TEXT NOT NULL,
		oracle_text TEXT,
		image_uris JSONB NOT NULL,
		set TEXT NOT NULL,
		set_name TEXT NOT NULL
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

// EnsureDecksTable creates the commander_decks table if it doesn't exist.
func EnsureCommanderDecksTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS proto_commander_decks (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		name TEXT NOT NULL,
    	description TEXT,
		commander TEXT NOT NULL,
		cards TEXT[]
	);`
	_, err := db.Exec(query)
	return err
}
