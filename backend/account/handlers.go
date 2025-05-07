package account

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func generateSessionToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := CreateUser(DB, &user); err != nil {
		http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if user.ID == 0 {
		http.Error(w, "Failed to retrieve user ID after registration", http.StatusInternalServerError)
		return
	}

	token := generateSessionToken()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	createdAt := time.Now()
	_, err := DB.Exec("INSERT INTO sessions (user_id, token, expires_at, created_at) VALUES ($1, $2, $3, $4)", user.ID, token, expiresAt, createdAt)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := AuthenticateUser(DB, creds.Email, creds.Password)
	if err != nil {
		fmt.Println("here 1")
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token := generateSessionToken()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	createdAt := time.Now()
	_, err = DB.Exec("INSERT INTO sessions (user_id, token, expires_at, created_at) VALUES ($1, $2, $3, $4)", user.ID, token, expiresAt, createdAt)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err == nil && cookie.Value != "" {
		_, _ = DB.Exec("DELETE FROM sessions WHERE token = $1", cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		fmt.Println("here 2")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var user User
	err = DB.QueryRow(`
		SELECT u.id, u.email FROM users u
		JOIN sessions s ON s.user_id = u.id
		WHERE s.token = $1
	`, cookie.Value).Scan(&user.ID, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error fetching user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"email": user.Email})
}

func CreateUser(db *sql.DB, user *User) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id`
	return db.QueryRow(query, user.Username, user.Email, string(hashed)).Scan(&user.ID)
}

func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	u := &User{}
	query := `SELECT id, username, email, password FROM users WHERE email = $1`
	err := db.QueryRow(query, email).Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func AuthenticateUser(db *sql.DB, email, password string) (*User, error) {
	user, err := GetUserByEmail(db, email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}
	return user, nil
}
