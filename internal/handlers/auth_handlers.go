package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/vijayvenkatj/http-servsers/helpers"
	"github.com/vijayvenkatj/http-servsers/internal/auth"
	"github.com/vijayvenkatj/http-servsers/internal/database"
	"github.com/vijayvenkatj/http-servsers/internal/models"
)


func Login(config *models.ApiConfig) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var requestBody struct {
            Email     string `json:"email"`
            Password  string `json:"password"`
            ExpiresIn int    `json:"expires_in_seconds,omitempty"`
        }

        if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
            helpers.RespondWithError(w, 400, "Invalid request body")
            return
        }

        userData, err := config.DbQueries.GetUserByEmail(r.Context(), requestBody.Email)
        if err != nil {
            helpers.RespondWithError(w, 404, "User not found")
            return
        }

        if err := auth.CheckPasswordHash(requestBody.Password, userData.HashedPassword); err != nil {
            helpers.RespondWithError(w, 401, "Unauthorized")
            return
        }

        expiresIn := requestBody.ExpiresIn
        if expiresIn == 0 {
            expiresIn = 3600 // default 1 hour
        }

        jwtToken, err := auth.MakeJWT(
            userData.ID,
            config.JWT_SECRET,
            time.Duration(expiresIn)*time.Second,
        )
        if err != nil {
            helpers.RespondWithError(w, 500, "Unable to generate auth token")
            return
        }

        refreshToken, err := auth.MakeRefreshToken()
        if err != nil {
            helpers.RespondWithError(w, 500, "Unable to generate refresh token")
            return
        }

        err = config.DbQueries.StoreRefreshToken(r.Context(), database.StoreRefreshTokenParams{
            UserID:   userData.ID,
            Token:    refreshToken,
            ExpiresAt: time.Now().Add(time.Hour * 60 * 24),
        })
        if err != nil {
            helpers.RespondWithError(w, 500, err.Error())
            return
        }

        user := models.User{
            ID:        userData.ID,
            Email:     userData.Email,
            CreatedAt: userData.CreatedAt.Time,
            UpdatedAt: userData.UpdatedAt.Time,
            JwtToken:  jwtToken,
            RefreshToken: refreshToken,
            IsChirpyRed: userData.IsChirpyRed,
        }
        

        helpers.RespondWithJSON(w, 200, user)
    }
}

func RefreshAccessToken(config *models.ApiConfig) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        refreshToken := auth.GetBearerToken(r.Header);
        if refreshToken == "" {
            helpers.RespondWithError(w,401,"Unauthorised!");
            return
        }

        token, err := config.DbQueries.GetRefreshToken(r.Context(),refreshToken);
        if err != nil || token.ExpiresAt.Before(time.Now()) || token.RevokedAt.Valid {
            helpers.RespondWithError(w,401,"Unauthorised!");
            return 
        }

        jwtToken, err := auth.MakeJWT(token.UserID, config.JWT_SECRET, 3600 * time.Second);
        if err != nil {
            helpers.RespondWithError(w,401,"Unauthorised!");
            return 
        }

        helpers.RespondWithJSON(w, 200, &struct{ Token string `json:"token"` }{Token: jwtToken })
    }
}

func RevokeRefreshToken(config *models.ApiConfig) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        refreshToken := auth.GetBearerToken(r.Header);
        if refreshToken == "" {
            helpers.RespondWithError(w,401,"Unauthorised!");
            return
        }

        err := config.DbQueries.RevokeRefreshToken(r.Context(), refreshToken);
        if err != nil {
            helpers.RespondWithError(w,400,"Unable to revoke refresh token!");
            return 
        }

        w.WriteHeader(204);
    }
}

func UpdateUser(config *models.ApiConfig) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        accessToken := auth.GetBearerToken(r.Header);
        if accessToken == "" {
            helpers.RespondWithError(w,401,"Unauthorised")
            return
        }

        userId, err := auth.ValidateJWT(accessToken,config.JWT_SECRET);
        if err != nil {
            helpers.RespondWithError(w,401,"Unauthorised")
            return
        }

        var requestBody struct {
            Email string `json:"email"`
            Password string `json:"password"`
        }

        err = json.NewDecoder(r.Body).Decode(&requestBody);
        if err != nil {
            helpers.RespondWithError(w,400,"Unable to decode body")
            return
        }

        hashed_password, err := auth.HashPassword(requestBody.Password);
        if err != nil {
            helpers.RespondWithError(w,400,"Unable to hash password")
            return
        }

        userData, err := config.DbQueries.UpdateUser(r.Context(), database.UpdateUserParams{
            Email: requestBody.Email,
            HashedPassword: hashed_password,
            ID: userId,
        })
        if err != nil {
            helpers.RespondWithError(w,400,"Unable to update details")
            return
        }

        var user models.User = models.User{
            ID: userData.ID,
            Email: userData.Email,
            IsChirpyRed: userData.IsChirpyRed,
            CreatedAt: userData.CreatedAt.Time,
            UpdatedAt: userData.UpdatedAt.Time,
        }

        helpers.RespondWithJSON(w,200,user);
    }
}