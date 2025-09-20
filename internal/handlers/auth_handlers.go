package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/vijayvenkatj/http-servsers/helpers"
	"github.com/vijayvenkatj/http-servsers/internal/auth"
	"github.com/vijayvenkatj/http-servsers/internal/models"
)


func Login(config *models.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var requestBody struct {
			Email	string 	`json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&requestBody);
		if err != nil {
			helpers.RespondWithError(w,500,"Unable to decode body!")
		}

		userData, err := config.DbQueries.GetUserByEmail(r.Context(), requestBody.Email);
		if err != nil {
			helpers.RespondWithError(w,400,"Unable to find user!")
		}

		err = auth.CheckPasswordHash(requestBody.Password,userData.HashedPassword)
		if err != nil {
			helpers.RespondWithError(w,401,"UnAuthorised!")
		}

		var user models.User = models.User{
			ID: userData.ID,
			Email: userData.Email,
			CreatedAt: userData.CreatedAt.Time,
			UpdatedAt: userData.UpdatedAt.Time,
		}

		helpers.RespondWithJSON(w,200,user)
	}
}