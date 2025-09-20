package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/vijayvenkatj/http-servsers/helpers"
	"github.com/vijayvenkatj/http-servsers/internal/auth"
	"github.com/vijayvenkatj/http-servsers/internal/database"
	"github.com/vijayvenkatj/http-servsers/internal/models"
)

func CreateUser(config *models.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type","application/json");

		var requestBody struct{
			Email	string	`json:"email"`
			Password string `json:"password"`
		};

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			helpers.RespondWithError(w,500,"Unable to decode body!");
		}

		hashed_password, err := auth.HashPassword(requestBody.Password);
		if err != nil {
			helpers.RespondWithError(w,400,"Error creating hash object!");
		}

		databaseUser, err := config.DbQueries.CreateUser(r.Context(), database.CreateUserParams{
			Email: requestBody.Email,
			HashedPassword: hashed_password,
		});
		if err != nil {
			helpers.RespondWithError(w,400,"Error creating User object!");
		}

		var user models.User = models.User{
			ID: databaseUser.ID,
			CreatedAt: databaseUser.CreatedAt.Time,
			UpdatedAt: databaseUser.UpdatedAt.Time,
			Email: databaseUser.Email,
		}

		helpers.RespondWithJSON(w,201,user);
	}
}


func GetUsers(config *models.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type","application/json");

		usersData, err := config.DbQueries.GetUsers(r.Context());
		if err != nil {
			helpers.RespondWithError(w,400,"Error getting all users!")
			return
		}

		var users []models.User
		for _, u := range usersData {
			users = append(users, models.User{
				ID:        u.ID,
				CreatedAt: u.CreatedAt.Time,
				UpdatedAt: u.UpdatedAt.Time,
				Email:     u.Email,
			})
		}

		helpers.RespondWithJSON(w, 200, users)
	}
}

