package main

import (
	
	"RolandQuest/internal/database"
	"RolandQuest/internal/endpoints"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/rs/cors"
)

func BeginHTTPRedirect() {
	redirectHandler := func (w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r,
			"https://" + r.Host + r.URL.String(),
			http.StatusMovedPermanently)
	}
	redirect_router := chi.NewMux()
	redirect_router.Use(middleware.Logger)
	redirect_router.Handle("/*", http.HandlerFunc(redirectHandler))
	go http.ListenAndServe(os.Getenv("LISTEN_ADDR_HTTP"), redirect_router)
}

func main() {
	
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	
	if err := VideoDatabase.Initialize(os.Getenv("DB_PATH")); err != nil {
		log.Fatal(err)
	}
	
	router := chi.NewMux()
	router.Use(middleware.Logger)
	router.Handle("/public/*", public())
	router.Get("/", endpoints.WrapError(endpoints.Home))
	
	router.Get("/videolist/*", endpoints.WrapError(endpoints.VideoList))
	router.Get("/player/*", endpoints.WrapError(endpoints.Player))
	router.Get("/video/*", endpoints.WrapError(endpoints.VideoServe))
	router.Get("/homebrewery/*", endpoints.WrapError(endpoints.HomebreweryServe))
	
	router.Get("/alt", endpoints.WrapError(endpoints.Alt))
	router.Get("/alt/get", endpoints.WrapError(endpoints.AltGet))
	
	router.Get("/test", endpoints.WrapError(endpoints.Test))
	router.Post("/test/upload", endpoints.WrapError(endpoints.TestUpload));
	
	corsHandler := cors.Default().Handler(router)
	
	// HTTP
	listenAddr := os.Getenv("LISTEN_ADDR_HTTP")
	slog.Info("HTTP Server started.", "listenAddr", listenAddr);
	log.Fatal(http.ListenAndServe(listenAddr, corsHandler))
	
	// HTTPS
	// listenAddr := os.Getenv("LISTEN_ADDR_HTTPS")
	// slog.Info("HTTP Server started.", "listenAddr", listenAddr);
	// BeginHTTPRedirect()
	// log.Fatal(http.ListenAndServeTLS(listenAddr, "test-serv-cert.pem", "test-priv-key.pem", corsHandler))
}
