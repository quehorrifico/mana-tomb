package middleware

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// DB is a global or injected reference to your database. If you're not using dependency injection, you'll need to set this during app startup.
var DB *sql.DB

func AuthMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil || cookie.Value == "" {
			fmt.Println("here 4")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := validateSessionToken(cookie.Value) // Replace with your actual logic
		if err != nil {
			fmt.Println("here 5")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateSessionToken(token string) (int, error) {
	var userID int
	var expiresAt time.Time

	err := DB.QueryRow(`SELECT user_id, expires_at FROM sessions WHERE token = $1`, token).Scan(&userID, &expiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("session not found")
		}
		return 0, err
	}

	if time.Now().After(expiresAt) {
		return 0, errors.New("session expired")
	}

	return userID, nil
}
