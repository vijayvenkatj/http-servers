package handlers

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/vijayvenkatj/http-servsers/helpers"
	"github.com/vijayvenkatj/http-servsers/internal/auth"
	"github.com/vijayvenkatj/http-servsers/internal/database"
	"github.com/vijayvenkatj/http-servsers/internal/models"
)

func CreateChirp(config *models.ApiConfig) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        var requestBody struct {
            Body string `json:"body"`
        }

        if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
            helpers.RespondWithError(w, 400, "Invalid request body")
            return
        }

        authToken := auth.GetBearerToken(r.Header)
        userID, err := auth.ValidateJWT(authToken, config.JWT_SECRET)
        if err != nil {
            helpers.RespondWithError(w, 401, "Unauthorized")
            return
        }

        chirpBody, err := helpers.ValidateChirp(requestBody.Body)
        if err != nil {
            helpers.RespondWithError(w, 400, "Invalid chirp")
            return
        }

        chirpData, err := config.DbQueries.CreateChirp(
            r.Context(),
            database.CreateChirpParams{
                UserID: userID,
                Body:   chirpBody,
            },
        )
        if err != nil {
            helpers.RespondWithError(w, 500, "Could not create chirp")
            return
        }

        chirp := models.Chirp{
            ID:        chirpData.ID,
            UserId:    chirpData.UserID,
            Body:      chirpData.Body,
            CreatedAt: chirpData.CreatedAt.Time,
            UpdatedAt: chirpData.UpdatedAt.Time,
        }

        helpers.RespondWithJSON(w, 201, chirp)
    }
}


func GetChirps(config *models.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type","application/json");

		var chirpsData []database.Chirp;

		author_id := r.URL.Query().Get("author_id");

		var authorUUID uuid.UUID
		var err error

		if author_id != "" {
			authorUUID, err = uuid.Parse(author_id)
			if err != nil {
				helpers.RespondWithError(w,400,"Error parsing author_id!")
				return
			}
			chirpsData, err = config.DbQueries.GetChirpsByAuthor(r.Context(), authorUUID)
			if err != nil {
				helpers.RespondWithError(w,400,"Error getting chirps by author!")
				return
			}
		} else {
			chirpsData, err = config.DbQueries.GetChirps(r.Context())
			if err != nil {
				helpers.RespondWithError(w,400,"Error getting all chirps!")
				return
			}
		}

		sorting := r.URL.Query().Get("sort");

		if sorting == "desc" {
			sort.Slice(chirpsData, func(i, j int) bool {
				return chirpsData[i].CreatedAt.Time.After(chirpsData[j].CreatedAt.Time)
			})
		} else {
			sort.Slice(chirpsData, func(i, j int) bool {
				return chirpsData[j].CreatedAt.Time.After(chirpsData[i].CreatedAt.Time)
			})
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

func DeleteChirp(config *models.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := r.PathValue("chirpId")
		chirpId, err := uuid.Parse(id)
		if err != nil {
			helpers.RespondWithError(w, 400, "Invalid chirp ID format")
			return
		}

		accessToken := auth.GetBearerToken(r.Header)
		if accessToken == "" {
			helpers.RespondWithError(w, 401, "Unauthorised")
			return
		}

		userId, err := auth.ValidateJWT(accessToken,config.JWT_SECRET);
		if err != nil {
			helpers.RespondWithError(w,403,"Forbidden")
			return 
		}

		chirp, err := config.DbQueries.GetChirpById(r.Context(),chirpId);
		if err != nil {
			helpers.RespondWithError(w,404,"Not found");
			return
		}

		if chirp.UserID != userId {
			helpers.RespondWithError(w,403,"Forbidden")
			return 
		}

		err = config.DbQueries.DeleteChirp(r.Context(), chirpId);
		if err != nil {
			helpers.RespondWithError(w,500,"Error deleting chirp!");
			return
		}

		w.WriteHeader(204);
	}
}