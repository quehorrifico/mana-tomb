package models

type Deck struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}
