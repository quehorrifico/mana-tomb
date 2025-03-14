package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// Wrapper for the response from Scryfall
type BulkDataResponse struct {
	Object  string     `json:"object"`
	HasMore bool       `json:"has_more"`
	Data    []BulkData `json:"data"`
}

// Individual Bulk Data struct
type BulkData struct {
	ID              string `json:"id"`
	URI             string `json:"uri"`
	Type            string `json:"type"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	DownloadURI     string `json:"download_uri"`
	UpdatedAt       string `json:"updated_at"`
	Size            int    `json:"size"`
	ContentType     string `json:"content_type"`
	ContentEncoding string `json:"content_encoding"`
}

func FetchAndStoreBulkData(db *sql.DB) {
	log.Println("Fetching Scryfall Bulk Data...")

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

	var bulkDataResponse BulkDataResponse // Use the wrapper struct

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
