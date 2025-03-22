package main

import (
	"log"
	"net/http"

	"github.com/quehorrifico/mana-tomb/account"
	"github.com/quehorrifico/mana-tomb/db"
	"github.com/quehorrifico/mana-tomb/handlers"
	"github.com/quehorrifico/mana-tomb/utils"
)

func main() {
	// 1) Connect to and open the final DB
	db.Connect()
	// Defer closing the final connection here, not inside Connect()
	defer db.GetDB().Close()

	// Link the DB to your handlers
	handlers.DB = db.GetDB()
	account.DB = db.GetDB()

	// 2) Ensure tables exist
	db.EnsureTables()

	// 3) Start your daily data fetch scheduler
	utils.StartScheduler(db.GetDB())

	// 4) Setup routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Mana Tomb API is running!"))
	})
	http.HandleFunc("/card/random", handlers.GetRandomCard)
	http.HandleFunc("/card/", handlers.GetCardByName)

	// Register account routes
	http.HandleFunc("/register", account.RegisterUser)
	http.HandleFunc("/login", account.LoginUser)

	// Additional routes for account, decks, etc. can go here

	log.Println("🚀 Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
