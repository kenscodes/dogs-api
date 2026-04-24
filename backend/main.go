package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

type Dog struct {
	ID        int64      `json:"id"`
	Breed     string     `json:"breed"`
	SubBreeds *string    `json:"subBreeds,omitempty"`
	Active    bool       `json:"active"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type DogRequest struct {
	Breed     string  `json:"breed"`
	SubBreeds *string `json:"subBreeds,omitempty"`
}

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

var db *sql.DB

func main() {
	var err error
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "../dogs.db"
	}
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create table
	createTable()

	// Initialize data if empty
	initializeData()

	// Setup router
	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/api/dogs", getDogs).Methods("GET")
	r.HandleFunc("/api/dogs/{id}", getDog).Methods("GET")
	r.HandleFunc("/api/dogs/breed/{breed}", getDogByBreed).Methods("GET")
	r.HandleFunc("/api/dogs/search", searchDogs).Methods("GET")
	r.HandleFunc("/api/dogs", createDog).Methods("POST")
	r.HandleFunc("/api/dogs/{id}", updateDog).Methods("PUT")
	r.HandleFunc("/api/dogs/{id}", deleteDog).Methods("DELETE")

	// Serve static files
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("../frontend")))

	// CORS handler
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS dogs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		breed TEXT NOT NULL UNIQUE,
		sub_breeds TEXT,
		active BOOLEAN DEFAULT 1,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

	// Add columns if they don't exist (for existing databases)
	_, err = db.Exec("ALTER TABLE dogs ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP")
	if err != nil {
		// Column might already exist, ignore error
	}
	_, err = db.Exec("ALTER TABLE dogs ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP")
	if err != nil {
		// Column might already exist, ignore error
	}
}

func initializeData() {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM dogs").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count == 0 {
		// Read dogs.json
		data, err := os.ReadFile("./dogs.json")
		if err != nil {
			log.Printf("Warning: Could not read dogs.json: %v", err)
			return
		}

		var dogsData map[string][]string
		err = json.Unmarshal(data, &dogsData)
		if err != nil {
			log.Printf("Warning: Could not parse dogs.json: %v", err)
			return
		}

		// Insert dogs
		for breed, subBreeds := range dogsData {
			var subBreedsStr *string
			if len(subBreeds) > 0 {
				s := strings.Join(subBreeds, ", ")
				subBreedsStr = &s
			}
			_, err := db.Exec(
				"INSERT INTO dogs (breed, sub_breeds, active) VALUES (?, ?, 1)",
				breed, subBreedsStr,
			)
			if err != nil {
				log.Printf("Error inserting breed %s: %v", breed, err)
			}
		}
		fmt.Printf("Initialized %d dog breeds\n", len(dogsData))
	}
}

func getDogs(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, breed, sub_breeds, active, created_at, updated_at FROM dogs WHERE active = 1 ORDER BY updated_at DESC")
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	dogs := []Dog{}
	for rows.Next() {
		var dog Dog
		err := rows.Scan(&dog.ID, &dog.Breed, &dog.SubBreeds, &dog.Active, &dog.CreatedAt, &dog.UpdatedAt)
		if err != nil {
			sendError(w, http.StatusInternalServerError, err.Error())
			return
		}
		dogs = append(dogs, dog)
	}

	sendResponse(w, http.StatusOK, "Retrieved all dogs", dogs)
}

func getDog(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var dog Dog
	err := db.QueryRow(
		"SELECT id, breed, sub_breeds, active, created_at, updated_at FROM dogs WHERE id = ? AND active = 1",
		id,
	).Scan(&dog.ID, &dog.Breed, &dog.SubBreeds, &dog.Active, &dog.CreatedAt, &dog.UpdatedAt)

	if err == sql.ErrNoRows {
		sendError(w, http.StatusNotFound, "Dog not found")
		return
	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusOK, "Retrieved dog", dog)
}

func getDogByBreed(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	breed := params["breed"]

	var dog Dog
	err := db.QueryRow(
		"SELECT id, breed, sub_breeds, active, created_at, updated_at FROM dogs WHERE breed = ? AND active = 1",
		breed,
	).Scan(&dog.ID, &dog.Breed, &dog.SubBreeds, &dog.Active, &dog.CreatedAt, &dog.UpdatedAt)

	if err == sql.ErrNoRows {
		sendError(w, http.StatusNotFound, "Dog not found")
		return
	}
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusOK, "Retrieved dog", dog)
}

func searchDogs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		getDogs(w, r)
		return
	}

	query = "%" + strings.ToLower(query) + "%"
	rows, err := db.Query(
		"SELECT id, breed, sub_breeds, active, created_at, updated_at FROM dogs WHERE active = 1 AND (LOWER(breed) LIKE ? OR LOWER(sub_breeds) LIKE ?) ORDER BY updated_at DESC",
		query, query,
	)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	dogs := []Dog{}
	for rows.Next() {
		var dog Dog
		err := rows.Scan(&dog.ID, &dog.Breed, &dog.SubBreeds, &dog.Active, &dog.CreatedAt, &dog.UpdatedAt)
		if err != nil {
			sendError(w, http.StatusInternalServerError, err.Error())
			return
		}
		dogs = append(dogs, dog)
	}

	sendResponse(w, http.StatusOK, "Search results", dogs)
}

func createDog(w http.ResponseWriter, r *http.Request) {
	var req DogRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Breed == "" {
		sendError(w, http.StatusBadRequest, "Breed name is required")
		return
	}

	// Check if breed already exists
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM dogs WHERE breed = ?", req.Breed).Scan(&count)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if count > 0 {
		sendError(w, http.StatusBadRequest, "Dog breed already exists")
		return
	}

	result, err := db.Exec(
		"INSERT INTO dogs (breed, sub_breeds, active) VALUES (?, ?, 1)",
		req.Breed, req.SubBreeds,
	)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	id, _ := result.LastInsertId()

	// Fetch the created dog with timestamps
	var dog Dog
	err = db.QueryRow(
		"SELECT id, breed, sub_breeds, active, created_at, updated_at FROM dogs WHERE id = ?",
		id,
	).Scan(&dog.ID, &dog.Breed, &dog.SubBreeds, &dog.Active, &dog.CreatedAt, &dog.UpdatedAt)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusCreated, fmt.Sprintf("Added %s", dog.Breed), dog)
}

func updateDog(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var req DogRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Breed == "" {
		sendError(w, http.StatusBadRequest, "Breed name is required")
		return
	}

	// Check if breed exists (for different breed)
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM dogs WHERE breed = ? AND id != ?", req.Breed, id).Scan(&count)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if count > 0 {
		sendError(w, http.StatusBadRequest, "Dog breed already exists")
		return
	}

	_, err = db.Exec(
		"UPDATE dogs SET breed = ?, sub_breeds = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		req.Breed, req.SubBreeds, id,
	)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Fetch the updated dog with timestamps
	var dog Dog
	err = db.QueryRow(
		"SELECT id, breed, sub_breeds, active, created_at, updated_at FROM dogs WHERE id = ?",
		id,
	).Scan(&dog.ID, &dog.Breed, &dog.SubBreeds, &dog.Active, &dog.CreatedAt, &dog.UpdatedAt)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusOK, fmt.Sprintf("Updated %s", dog.Breed), dog)
}

func deleteDog(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	// Fetch breed name before deleting
	var breed string
	err := db.QueryRow("SELECT breed FROM dogs WHERE id = ?", id).Scan(&breed)
	if err != nil {
		sendError(w, http.StatusNotFound, "Dog not found")
		return
	}

	_, err = db.Exec("UPDATE dogs SET active = 0 WHERE id = ?", id)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusOK, fmt.Sprintf("Deleted %s", breed), nil)
}

func sendResponse(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func sendError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ApiResponse{
		Success: false,
		Message: message,
	})
}

func parseID(s string) int64 {
	var id int64
	fmt.Sscanf(s, "%d", &id)
	return id
}
