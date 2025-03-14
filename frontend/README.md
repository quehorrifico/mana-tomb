# Frontend
The frontend is built with React + TypeScript and connects to the Go backend.

Setup Instructions
	1.	Install dependencies
        npm install
    2.	Run the development server
        npm run dev
    3.	Build the project (for production)
        npm run build

Project Structure
frontend/
│── src/
│   ├── components/    # UI Components
│   ├── pages/         # Page views
│   ├── api/           # API calls
│   ├── App.tsx        # Main entry point
│── public/            # Static assets
│── package.json       # Dependencies
│── README.md          # Frontend documentation

Fetching a Random Card

The frontend fetches a random card from the backend using the /api/random-card endpoint.

Known Issues & Fixes
	•	Frontend does not display data
        •	Ensure the backend is running.
        •	Check the browser console for errors.