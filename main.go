package main

import (
	"Chirpy/internal/database"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits int
	db             *database.Queries
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("postgres://shihong:@localhost:5432/gator?sslmode=disable")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	dbQueries := database.New(db)
	apiCfg := apiConfig{
		fileserverHits: 0,
		db:             dbQueries,
	}
	mux := http.NewServeMux()
	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))

	})
	mux.HandleFunc("POST /api/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(405)
	})
	mux.HandleFunc("DELETE /api/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(405)
	})
	mux.Handle("GET /admin/metrics", apiCfg.hitCounter())
	mux.HandleFunc("POST /admin/metrics", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(405)
	})
	mux.HandleFunc("DELETE /admin/metrics", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(405)
	})
	mux.Handle("/api/reset", apiCfg.resetCounter())
	mux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
		type params struct {
			Body string `json:"body"`
		}
		type valid struct {
			Valid bool `json:"valid"`
		}
		decoder := json.NewDecoder(r.Body)
		param := params{}
		err := decoder.Decode(&param)
		if err != nil {
			msg := "Something went wrong"
			respondWithError(w, 400, msg)
			return
		}
		if len(string(param.Body)) > 140 {
			msg := "Chirp is too long"
			respondWithError(w, 400, msg)
			return
		}
		stringList := strings.Split(string(param.Body), " ")
		var resString []string
		for _, v := range stringList {
			if strings.Contains(strings.ToLower(v), "kerfuffle") || strings.Contains(strings.ToLower(v), "sharbert") || strings.Contains(strings.ToLower(v), "fornax") {
				resString = append(resString, "****")
			} else {
				resString = append(resString, v)
			}
		}
		result := strings.Join(resString, " ")
		type cleanBody struct {
			Cleaned_body string `json:"cleaned_body"`
		}
		cleaned := cleanBody{
			Cleaned_body: result,
		}
		respondWithJSON(w, 200, cleaned)
	})
	err = http.ListenAndServe(s.Addr, s.Handler)
	if err != nil {
		fmt.Println("Ecounter Error during Server Initiation:", err)
		return
	}
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
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

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type rtnErr struct {
		Error string `json:"error"`
	}
	decodeErr := rtnErr{
		Error: msg,
	}
	dat, err := json.Marshal(decodeErr)
	if err != nil {
		fmt.Println(err)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
