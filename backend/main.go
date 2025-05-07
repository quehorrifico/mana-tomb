package main

import (
	"log"
	"net/http"

	"github.com/quehorrifico/mana-tomb/backend/account"
	"github.com/quehorrifico/mana-tomb/backend/cards"
	"github.com/quehorrifico/mana-tomb/backend/db"
	"github.com/quehorrifico/mana-tomb/backend/decks"
	"github.com/quehorrifico/mana-tomb/backend/middleware"
	"github.com/quehorrifico/mana-tomb/backend/utils"
)

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Only allow frontend origin
		if origin == "http://localhost:3000" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		}

		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cookie")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func registerRoutes(mux *http.ServeMux) {
	// Health check
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Mana Tomb API is running!"))
	})

	// Card endpoints (Public)
	mux.Handle("/card/random", withCORS(http.HandlerFunc(cards.GetRandomCard)))
	mux.Handle("/card/", withCORS(http.HandlerFunc(cards.GetCardByName)))

	// Deck endpoints (Protected)
	mux.Handle("/decks", withCORS(middleware.AuthMiddleware(http.HandlerFunc(decks.GetDecksByUser))))
	mux.Handle("/decks/create", withCORS(middleware.AuthMiddleware(http.HandlerFunc(decks.CreateDeck))))
	mux.Handle("/decks/", withCORS(middleware.AuthMiddleware(http.HandlerFunc(decks.GetDeckByID))))
	mux.Handle("/decks/update/", withCORS(middleware.AuthMiddleware(http.HandlerFunc(decks.UpdateDeck))))
	mux.Handle("/decks/delete/", withCORS(middleware.AuthMiddleware(http.HandlerFunc(decks.DeleteDeck))))

	// Auth endpoints (Public)
	mux.Handle("/register", withCORS(http.HandlerFunc(account.RegisterUser)))
	mux.Handle("/login", withCORS(http.HandlerFunc(account.LoginUser)))
	mux.Handle("/logout", withCORS(http.HandlerFunc(account.LogoutUser)))
	mux.Handle("/me", withCORS(http.HandlerFunc(account.GetCurrentUser)))
}

func main() {
	// 1) Connect to and open the final DB
	db.Connect()
	defer db.GetDB().Close()

	// 2) Inject DB references into packages
	cards.DB = db.GetDB()
	account.DB = db.GetDB()
	decks.DB = db.GetDB()
	middleware.DB = db.GetDB()

	// 3) Initialize database schema and start daily card fetch job
	utils.StartScheduler(db.GetDB())

	// 4) Setup HTTP routes and start the server
	mux := http.NewServeMux()
	registerRoutes(mux)

	log.Println("🚀 Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
