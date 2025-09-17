package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
)


type apiConfig struct {
	fileServerHits atomic.Int32
	plainTextReqs  atomic.Int32
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



func main() {

	serveMux := http.NewServeMux();

	config := &apiConfig{
		fileServerHits: atomic.Int32{},
	}

	serveMux.Handle("/app/", config.HandleMetricMiddleware(http.StripPrefix("/app",http.FileServer(http.Dir("./")))));

	serveMux.Handle("POST /reset", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config.fileServerHits.Store(0)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("metrics reset"))
	}))

	serveMux.Handle("GET /metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Hits: %d", config.fileServerHits.Load())))
	}))

	serveMux.Handle("GET /healthz",http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type","text/plain; charset=utf-8")
		w.WriteHeader(200)
		body := "OK"
		w.Write([]byte(body));
	}));


	server := &http.Server{
		Addr: ":8080",
		Handler: serveMux,
	}

	err := server.ListenAndServe();
	if err != nil {
		return 
	}

}