package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/vijayvenkatj/http-servsers/helpers"
	"github.com/vijayvenkatj/http-servsers/internal/database"
	"github.com/vijayvenkatj/http-servsers/internal/models"
)


type apiConfig struct {
	fileServerHits atomic.Int32
	plainTextReqs  atomic.Int32
	dbQueries	*database.Queries
}

func (config *apiConfig) HandleMetricMiddleware(next http.Handler) http.Handler {
	wrapperReq := config.Middleware(next);
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control","no-cache")
		config.fileServerHits.Add(1);
        wrapperReq.ServeHTTP(w, r)
    })
}

func (config *apiConfig) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("Content-Type");

		if strings.HasPrefix(header, "text/plain") {
			config.plainTextReqs.Add(1);
		}

        next.ServeHTTP(w, r)
    })
}


func init() {
	godotenv.Load()
}


func main() {

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error : %s", err);
		return
	}

	dbQueries := database.New(db);

	serveMux := http.NewServeMux();

	config := &apiConfig{
		fileServerHits: atomic.Int32{},
		dbQueries: dbQueries,
	}

	serveMux.Handle("/app/", config.HandleMetricMiddleware(http.StripPrefix("/app",http.FileServer(http.Dir("./")))));
	serveMux.Handle("GET /api/healthz",http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type","text/plain; charset=utf-8")
		w.WriteHeader(200)
		body := "OK"
		w.Write([]byte(body));
	}));


	serveMux.Handle("GET /admin/metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type","text/html")

		w.Write([]byte(fmt.Sprintf(`
		<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
		</html>
		`, config.fileServerHits.Load())))


	}))

	serveMux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Content-Type","application/json");

		var requestBody struct{
			Email	string	`json:"email"`
		};

		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			helpers.RespondWithError(w,500,"Unable to decode body!");
		}

		databaseUser, err := config.dbQueries.CreateUser(r.Context(), requestBody.Email);
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
	})

	serveMux.HandleFunc("POST /api/chirps", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		chirpData, err := config.dbQueries.CreateChirp(r.Context(),
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
	}))

	serveMux.HandleFunc("GET /api/users", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type","application/json");

		usersData, err := config.dbQueries.GetUsers(r.Context());
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
	}))

	serveMux.HandleFunc("GET /api/chirps", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type","application/json");

		chirpsData, err := config.dbQueries.GetChirps(r.Context());
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
	}))

	serveMux.HandleFunc("GET /api/chirps/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			helpers.RespondWithError(w, 400, "Invalid chirp ID format")
			return
		}

		chirp, err := config.dbQueries.GetChirpById(r.Context(), id)
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
	})

	serveMux.Handle("POST /admin/reset", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config.fileServerHits.Store(0)

		platform := os.Getenv("PLATFORM");
		if platform == "dev" {
			err = config.dbQueries.DeleteAllUsers(r.Context());
			if err != nil {
				helpers.RespondWithError(w,400,"Error deleting all users!")
				return
			}

			helpers.RespondWithJSON(w,200,"Users deleted successfully!");
		}

		helpers.RespondWithJSON(w,403,"Forbidden!");
	}))

	server := &http.Server{
		Addr: ":8080",
		Handler: serveMux,
	}

	err = server.ListenAndServe();
	if err != nil {
		return 
	}

}