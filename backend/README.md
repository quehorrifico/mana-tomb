# Backend
The backend is written in Go and serves an API for querying Magic: The Gathering cards.

Environment Variables

Create a .env file in /backend/ with:
    DB_USER=mana_user
    DB_PASSWORD=mana_pass
    DB_NAME=mana_tomb
    DB_HOST=localhost
    DB_PORT=5432
    SCRYFALL_BULK_URL=https://api.scryfall.com/bulk-data

How to Run the Backend
    go run main.go

Database Management
	•	Connect to PostgreSQL inside Docker
            docker exec -it mana-postgres psql -U postgres -d mana_tomb
    •	List all tables
            \dt
	•	View table data
            SELECT * FROM cards LIMIT 5;
	•	Drop a table
            DROP TABLE cards;

API Endpoints
    Endpoint	Method	Description
    /api/random-card	GET	Returns a random MTG card
    /api/cards?name=	GET	Search cards by name
    /api/bulk-update	POST	Manually trigger a Scryfall update

## 🏗 Project Structure
backend/
├── account/             # User registration/login logic
│   ├── handlers.go
│   └── models.go
├── cards/               # HTTP handlers for card search and random card
│   └── handlers.go
├── decks/               # Deck builder logic (WIP)
│   ├── handlers.go
│   └── models.go
├── db/                  # DB connection + initialization
│   └── connection.go
├── models/              # Shared DB models (PostgreSQL schemas)
│   ├── bulk_data.go
│   ├── deck.go
│   ├── deck_card.go
│   ├── oracle_card.go
│   └── unique_artwork.go
├── middleware/          # Middleware like CORS
│   └── cors.go
├── utils/               # Shared tools and scheduled jobs
│   ├── parser.go
│   ├── scheduler.go
│   └── setup.go
├── main.go              # Server startup and route registration
└── .env                 # Environment config (Postgres URL, etc.)

