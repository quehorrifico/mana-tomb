package main

import (
	"log"
	"net/http"

	"github.com/quehorrifico/mana-tomb/account"
	"github.com/quehorrifico/mana-tomb/db"
	"github.com/quehorrifico/mana-tomb/handlers"
	"github.com/quehorrifico/mana-tomb/utils"
)

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.EnableCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func main() {
	// 1) Connect to and open the final DB
	db.Connect()
	defer db.GetDB().Close()

	// Link the DB to your handlers
	handlers.DB = db.GetDB()
	account.DB = db.GetDB()

	// 2) Ensure tables exist
	db.EnsureTables()

	// 3) Start your daily data fetch scheduler
	utils.StartScheduler(db.GetDB())

	// 4) Setup ServeMux and routes
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Mana Tomb API is running!"))
	})
	mux.HandleFunc("/card/random", withCORS(handlers.GetRandomCard))
	mux.HandleFunc("/card/", withCORS(handlers.GetCardByName))
	mux.HandleFunc("/register", withCORS(account.RegisterUser))
	mux.HandleFunc("/login", withCORS(account.LoginUser))
	mux.HandleFunc("/logout", withCORS(account.LogoutUser))
	mux.HandleFunc("/me", withCORS(account.GetCurrentUser))

	log.Println("🚀 Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
