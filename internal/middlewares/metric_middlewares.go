package middlewares

import (
	"net/http"
	"strings"

	"github.com/vijayvenkatj/http-servsers/internal/models"
)


func HandleMetricMiddleware(next http.Handler, config *models.ApiConfig) http.Handler {
	wrapperReq := Middleware(next,config);
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control","no-cache")
		config.FileServerHits.Add(1);
        wrapperReq.ServeHTTP(w, r)
    })
}

func Middleware(next http.Handler, config *models.ApiConfig) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("Content-Type");

		if strings.HasPrefix(header, "text/plain") {
			config.PlainTextReqs.Add(1);
		}

        next.ServeHTTP(w, r)
    })
}