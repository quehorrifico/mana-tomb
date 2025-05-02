package main

import (
	"log"
	"net/http"

	"github.com/quehorrifico/mana-tomb/backend/account"
	"github.com/quehorrifico/mana-tomb/backend/cards"
	"github.com/quehorrifico/mana-tomb/backend/db"
	"github.com/quehorrifico/mana-tomb/backend/middleware"
	"github.com/quehorrifico/mana-tomb/backend/utils"
)

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		middleware.EnableCORS(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

func registerRoutes(mux *http.ServeMux) {
	// Health check
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Mana Tomb API is running!"))
	})

	// Card endpoints
	mux.HandleFunc("/card/random", withCORS(cards.GetRandomCard))
	mux.HandleFunc("/card/", withCORS(cards.GetCardByName))

	// Auth endpoints
	mux.HandleFunc("/register", withCORS(account.RegisterUser))
	mux.HandleFunc("/login", withCORS(account.LoginUser))
	mux.HandleFunc("/logout", withCORS(account.LogoutUser))
	mux.HandleFunc("/me", withCORS(account.GetCurrentUser))
}

func main() {
	// 1) Connect to and open the final DB
	db.Connect()
	defer db.GetDB().Close()

	// 2) Inject DB references into packages
	cards.DB = db.GetDB()
	account.DB = db.GetDB()

	// 3) Initialize database schema and start daily card fetch job
	utils.StartScheduler(db.GetDB())

	// 4) Setup HTTP routes and start the server
	mux := http.NewServeMux()
	registerRoutes(mux)

	log.Println("🚀 Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
