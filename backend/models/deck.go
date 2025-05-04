package models

type ProtoCommanderDeck struct {
	UserID      int      `json:"user_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Commander   string   `json:"commander"`
	Cards       []string `json:"cards"`
}
