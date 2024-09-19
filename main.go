package main

import (
	"net/http"
	"fmt"
)

type apiConfig struct {
	fileserverHits int
}

func main(){
	apiCfg := apiConfig{
		fileserverHits: 0,
	}
	mux := http.NewServeMux()
	s := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/",http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, req *http.Request){
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))

	})
	mux.HandleFunc("POST /healthz", func(w http.ResponseWriter, req *http.Request){
		w.WriteHeader(405)
	})
	mux.HandleFunc("DELETE /healthz", func(w http.ResponseWriter, req *http.Request){
		w.WriteHeader(405)
	})
	mux.Handle("GET /metrics", apiCfg.hitCounter())
	mux.HandleFunc("POST /metrics", func(w http.ResponseWriter, req *http.Request){
		w.WriteHeader(405)
	})
	mux.HandleFunc("DELETE /metrics", func(w http.ResponseWriter, req *http.Request){
		w.WriteHeader(405)
	})
	mux.Handle("/reset", apiCfg.resetCounter())
	err := http.ListenAndServe(s.Addr, s.Handler)
	if err != nil{
		fmt.Println("Ecounter Error during Server Initiation:", err)
		return
	}
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits ++
		next.ServeHTTP(w,r)
	})
}

func (cfg *apiConfig) hitCounter() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := fmt.Sprintf("Hits: %d", cfg.fileserverHits)
		w.Write([]byte(resp))
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	})
}

func (cfg *apiConfig) resetCounter() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits = 0
		w.WriteHeader(200)
	})
}
