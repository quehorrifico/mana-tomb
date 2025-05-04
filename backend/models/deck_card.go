package models

type DeckCard struct {
	ID       int    `json:"id"`
	DeckID   int    `json:"deck_id"`
	CardID   string `json:"card_id"` // Foreign key to OracleCard/UniqueArtworkCard
	CardName string `json:"card_name"`
	Quantity int    `json:"quantity"`
}
