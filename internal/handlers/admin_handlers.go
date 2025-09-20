package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/vijayvenkatj/http-servsers/helpers"
	"github.com/vijayvenkatj/http-servsers/internal/models"
)


func GetMetrics(config *models.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type","text/html")

		w.Write([]byte(fmt.Sprintf(`
		<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
		</html>
		`, config.FileServerHits.Load())))
	}
}

func Reset(config *models.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		config.FileServerHits.Store(0)

		platform := os.Getenv("PLATFORM");
		if platform == "dev" {
			err := config.DbQueries.DeleteAllUsers(r.Context());
			if err != nil {
				helpers.RespondWithError(w,400,"Error deleting all users!")
				return
			}

			helpers.RespondWithJSON(w,200,"Users deleted successfully!");
		}

		helpers.RespondWithJSON(w,403,"Forbidden!");
	}
}

func HealthZ(config *models.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type","text/plain; charset=utf-8")
		w.WriteHeader(200)
		body := "OK"
		w.Write([]byte(body));
	}
}