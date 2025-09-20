package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/vijayvenkatj/http-servsers/internal/database"
	"github.com/vijayvenkatj/http-servsers/internal/handlers"
	"github.com/vijayvenkatj/http-servsers/internal/middlewares"
	"github.com/vijayvenkatj/http-servsers/internal/models"
)



var Config *models.ApiConfig

func init() {
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error : %s", err);
		return
	}

	DbQueries := database.New(db);
	Config = &models.ApiConfig{
		FileServerHits: atomic.Int32{},
		DbQueries: DbQueries,
	}
}


func main() {

	serveMux := http.NewServeMux();

	serveMux.Handle("/app/", middlewares.HandleMetricMiddleware(http.StripPrefix("/app",http.FileServer(http.Dir("./"))),Config));

	serveMux.HandleFunc("POST /api/users", handlers.CreateUser(Config));
	serveMux.HandleFunc("GET /api/users", handlers.GetUsers(Config));

	serveMux.HandleFunc("POST /api/login", handlers.Login(Config));

	serveMux.HandleFunc("POST /api/chirps", handlers.CreateChirp(Config));
	serveMux.HandleFunc("GET /api/chirps", handlers.GetChirps(Config));
	serveMux.HandleFunc("GET /api/chirps/{id}", handlers.GetChirpById(Config));

	serveMux.Handle("GET /admin/metrics", handlers.GetMetrics(Config))
	serveMux.Handle("POST /admin/reset", handlers.Reset(Config))
	serveMux.Handle("GET /api/healthz",handlers.HealthZ(Config));

	server := &http.Server{
		Addr: ":8080",
		Handler: serveMux,
	}

	err := server.ListenAndServe();
	if err != nil {
		return 
	}

}