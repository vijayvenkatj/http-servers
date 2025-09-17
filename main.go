package main


import (
	"net/http"
)

type healthHandler struct {}

func (healthHandler) ServeHTTP(w http.ResponseWriter,r *http.Request){

	w.Header().Set("Content-Type","text/plain; charset=utf-8")

	w.WriteHeader(200)

	body := "OK"
	w.Write([]byte(body));
}


func main() {

	serveMux := http.NewServeMux();

	serveMux.Handle("/app/", http.StripPrefix("/app",http.FileServer(http.Dir("./"))));
	serveMux.Handle("/healthz",healthHandler{});

	server := &http.Server{
		Addr: ":8080",
		Handler: serveMux,
	}

	err := server.ListenAndServe();
	if err != nil {
		return 
	}

}