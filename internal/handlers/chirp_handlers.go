package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/vijayvenkatj/http-servsers/helpers"
	"github.com/vijayvenkatj/http-servsers/internal/database"
	"github.com/vijayvenkatj/http-servsers/internal/models"
)


func CreateChirp(config *models.ApiConfig) http.HandlerFunc {
	return (func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type","application/json");

		var requestBody struct {
			Body	string `json:"body"`
			UserId	uuid.UUID `json:"user_id"`
		};

		err := json.NewDecoder(r.Body).Decode(&requestBody);
		if err != nil {
			helpers.RespondWithError(w,500,"Unable to decode body!");
			return
		}

		chirpBody, err := helpers.ValidateChirp(requestBody.Body);
		if err != nil {
			helpers.RespondWithError(w,400,"Unable to validate chirp!");
			return
		}

		chirpData, err := config.DbQueries.CreateChirp(r.Context(),
			database.CreateChirpParams{
				UserID: requestBody.UserId,
				Body: chirpBody,
			},
		)
		if err != nil {
			helpers.RespondWithError(w,400,"Unable to create chirp!");
			return
		}

		var chirp models.Chirp = models.Chirp {
			ID: chirpData.ID,
			UserId: chirpData.UserID,
			Body: chirpData.Body,
			CreatedAt: chirpData.CreatedAt.Time,
			UpdatedAt: chirpData.UpdatedAt.Time,
		}

		helpers.RespondWithJSON(w,201,chirp);
	})
}


func GetChirps(config *models.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type","application/json");

		chirpsData, err := config.DbQueries.GetChirps(r.Context());
		if err != nil {
			helpers.RespondWithError(w,400,"Error getting all users!")
			return
		}

		var chirps []models.Chirp
		for _, u := range chirpsData {
			chirps = append(chirps, models.Chirp{
				ID:        u.ID,
				CreatedAt: u.CreatedAt.Time,
				UpdatedAt: u.UpdatedAt.Time,
				Body:      u.Body,
				UserId:    u.UserID,
			})
		}

		helpers.RespondWithJSON(w, 200, chirps)
	}
}

func GetChirpById(config *models.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			helpers.RespondWithError(w, 400, "Invalid chirp ID format")
			return
		}

		chirp, err := config.DbQueries.GetChirpById(r.Context(), id)
		if err != nil {
			helpers.RespondWithError(w, 404, "Chirp not found")
			return
		}

		helpers.RespondWithJSON(w, 200, models.Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt.Time,
			UpdatedAt: chirp.UpdatedAt.Time,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		})
	}
}