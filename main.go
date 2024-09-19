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
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request){
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))

	})
	mux.HandleFunc("POST /api/healthz", func(w http.ResponseWriter, req *http.Request){
		w.WriteHeader(405)
	})
	mux.HandleFunc("DELETE /api/healthz", func(w http.ResponseWriter, req *http.Request){
		w.WriteHeader(405)
	})
	mux.Handle("GET /admin/metrics", apiCfg.hitCounter())
	mux.HandleFunc("POST /admin/metrics", func(w http.ResponseWriter, req *http.Request){
		w.WriteHeader(405)
	})
	mux.HandleFunc("DELETE /admin/metrics", func(w http.ResponseWriter, req *http.Request){
		w.WriteHeader(405)
	})
	mux.Handle("/api/reset", apiCfg.resetCounter())
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
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		htmlTemplate := fmt.Sprintf(`<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>`, cfg.fileserverHits)
		w.Write([]byte(htmlTemplate))
	})
}

func (cfg *apiConfig) resetCounter() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits = 0
		w.WriteHeader(200)
	})
}
