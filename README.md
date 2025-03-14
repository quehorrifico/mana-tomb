# Mana Tomb
Mana Tomb is a web application that provides a searchable database of Magic: The Gathering cards using the Scryfall API and a self-hosted PostgreSQL database.

# Project Structure
mana-tomb/
│── backend/            # Go backend API
│── frontend/           # React + TypeScript frontend
│── docker/             # Docker configurations
│── README.md           # General project overview
│── init-db.sh          # Database initialization script
│── .env                # Environment variables (not committed)

# Getting Started
	1.	Ensure you have these dependencies installed:
	•	Go (for the backend)
	•	Node.js & npm (for the frontend)
	•	Docker (for database management)
	•	PostgreSQL (if running locally)

	2.	Start the database with Docker
        cd docker
        docker-compose up -d

    3.	Start the backend server
        cd backend
        go run main.go
    
    4.	Start the frontend
    cd frontend
    npm install  # First time only
    npm run dev  # Start frontend

    5.	Check if everything is working
        •	Backend: Open http://localhost:8080/api/random-card to check if the API is returning a random card.
        •	Frontend: Open http://localhost:3000 to see the app UI.
