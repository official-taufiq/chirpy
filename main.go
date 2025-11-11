package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/official-taufiq/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits *atomic.Int32
	db             *database.Queries
	platform       string
	tokenSecret    string
	apikey         string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	defer dbConn.Close()
	dbQueries := database.New(dbConn)
	err = dbConn.Ping()
	if err != nil {
		panic(err)
	}
	Platform := os.Getenv("PLATFORM")
	Secret := os.Getenv("JWT_SECRET")
	ApiKey := os.Getenv("POLKA_KEY")
	apiCfg := apiConfig{
		fileserverHits: &atomic.Int32{},
		db:             dbQueries,
		platform:       Platform,
		tokenSecret:    Secret,
		apikey:         ApiKey,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset) //Change to POST before final submission
	mux.HandleFunc("POST /api/chirps", apiCfg.chirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.AllChirps)
	mux.HandleFunc("POST /api/chirps/{chirpID}", apiCfg.OneChirp)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", apiCfg.login)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.tokenRevoke)
	mux.HandleFunc("PUT /api/users", apiCfg.updateUser)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.deleteChirp)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.chirpyRed)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())

}
