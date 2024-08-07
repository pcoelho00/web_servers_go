package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/pcoelho00/server_go/database"
	"github.com/pcoelho00/server_go/handlers"
)

func main() {

	const port = "8080"
	const root = "."
	const templates = root + "/templates"

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	PolkaKey := os.Getenv("POLKA_KEY")

	_, err = os.Stat("database.json")
	if !os.IsNotExist(err) {
		err := os.Remove("database.json")
		if err != nil {
			log.Fatal("Couldn't delete the database")
		}
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal("Can't connect with the Database")
	}

	apiCfg := handlers.ApiConfig{
		FileserverHits: 0,
		DB:             db,
		JwtSecret:      jwtSecret,
		PolkaKey:       PolkaKey,
	}

	mux := http.NewServeMux()
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(templates)))
	mux.Handle("/app/*", apiCfg.MiddlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", handlers.HealthsResponseHandler)

	mux.HandleFunc("GET /api/reset", apiCfg.ResetStatsHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.MetricsHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.PostChirpsHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.GetChirpsMsgHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.GetChirpHandler)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.DeleteChirpHandler)

	mux.HandleFunc("POST /api/users", apiCfg.PostUserHandler)
	mux.HandleFunc("GET /api/users", apiCfg.GetUsersHandler)
	mux.HandleFunc("PUT /api/users", apiCfg.PutLoginUserHandler)
	mux.HandleFunc("GET /api/users/{userID}", apiCfg.GetUserHandler)
	mux.HandleFunc("POST /api/login", apiCfg.PostLoginHandler)
	mux.HandleFunc("POST /api/refresh", apiCfg.RefreshTokenHandler)
	mux.HandleFunc("POST /api/revoke", apiCfg.RevokeRefreshTokenHandler)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.ChirpyRedHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Server started at port: %s\n", port)
	log.Fatal(server.ListenAndServe())

}
