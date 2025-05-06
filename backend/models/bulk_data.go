package models

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
