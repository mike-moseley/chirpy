package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/mike-moseley/chirpy/internal/auth"
	"github.com/mike-moseley/chirpy/internal/database"
)

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("OK"))
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	metricsStr := fmt.Sprintf(
		`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())

	w.Write([]byte(metricsStr))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Forbidden")
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	cfg.fileserverHits.Store(0)
	cfg.db.DeleteAllUsers(req.Context())
	resetStr := "File server hits reset to 0\nUsers deleted\n"

	w.Write([]byte(resetStr))
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type responseError struct {
		Error string `json:"error"`
	}
	respBody := responseError{
		Error: msg,
	}
	respondWithJSON(w, code, respBody)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(jsonData)
	if err != nil {
		log.Printf("Error writing response: %s", err)
	}
}

func replaceProfane(body string) string {
	cleanedBody := body
	var profane = []string{"kerfuffle", "sharbert", "fornax"}
	for s := range strings.SplitSeq(body, " ") {
		sLower := strings.ToLower(s)
		for _, p := range profane {
			if sLower == p {
				cleanedBody = strings.ReplaceAll(cleanedBody, s, "****")
			}
		}
	}
	return cleanedBody
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errMsg := fmt.Sprintf("Error decoding user: %s", err)
		respondWithError(w, 500, errMsg)
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		errMsg := fmt.Sprintf("Error decoding user: %s", err)
		respondWithError(w, 500, errMsg)
		return
	}
	dbUser, err := cfg.db.CreateUser(req.Context(), database.CreateUserParams{Email: params.Email, HashedPassword: hashedPassword})
	if err != nil {
		errMsg := fmt.Sprintf("Error creating user: %s", err)
		respondWithError(w, 500, errMsg)
		return
	}

	user := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	respondWithJSON(w, 201, user)
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body   string        `json:"body"`
		UserID uuid.NullUUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errMsg := fmt.Sprintf("Error decoding chirp: %s", err)
		respondWithError(w, 500, errMsg)
		return
	}
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	cleanedBody := replaceProfane(params.Body)
	dbChirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{Body: cleanedBody, UserID: params.UserID})
	if err != nil {
		errMsg := fmt.Sprintf("Error creating chirp: %s", err)
		respondWithError(w, 500, errMsg)
		return
	}
	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJSON(w, 201, chirp)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	dbChirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		errMsg := fmt.Sprintf("Error getting chirps: %s", err)
		respondWithError(w, 500, errMsg)
		return
	}
	chirps := make([]Chirp, len(dbChirps))
	for i, dbChirp := range dbChirps {
		chirps[i] = Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}
	}
	respondWithJSON(w, 200, chirps)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, req *http.Request) {
	chirpID := req.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		errMsg := fmt.Sprintf("Error converting string to UUID: %s", err)
		respondWithError(w, 500, errMsg)
		return
	}
	dbChirp, err := cfg.db.GetChirpByID(req.Context(), chirpUUID)
	if err != nil {
		errMsg := fmt.Sprintf("Error getting chirps: %s", err)
		respondWithError(w, 404, errMsg)
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJSON(w, 200, chirp)
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		errMsg := fmt.Sprintf("Error decoding username and password: %s", err)
		respondWithError(w, 500, errMsg)
		return
	}
	dbHashedPass, err := cfg.db.GetHashedPass(req.Context(), params.Email)
	if err != nil {
		errMsg := fmt.Sprintf("Error decoding username and password: %s", err)
		respondWithError(w, 500, errMsg)
		return
	}
	result, err := auth.CheckPasswordHash(params.Password, dbHashedPass)
	if !result || err != nil {
		errMsg := "Incorrect email or password"
		respondWithError(w, 401, errMsg)
		return
	}
	userJson, err := cfg.db.GetUser(req.Context(), params.Email)
	if err != nil {
		errMsg := "Incorrect email or password"
		respondWithError(w, 401, errMsg)
		return
	}
	user := User{
		ID:        userJson.ID,
		CreatedAt: userJson.CreatedAt,
		UpdatedAt: userJson.UpdatedAt,
		Email:     userJson.Email,
	}
	respondWithJSON(w, 200, user)
}
