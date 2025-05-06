package models

type ProtoCommanderDeck struct {
	DeckID      int      `json:"deck_id"`
	UserID      int      `json:"user_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Commander   string   `json:"commander"`
	Cards       []string `json:"cards"`
}
