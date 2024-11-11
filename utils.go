package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	isProduction bool
	addr         string
	staticPath   string
	indexPath    string
	Port         string
}

func loadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("błąd przy ładowaniu zmiennych środowiskowych")
	}

	isProduction := os.Getenv("MODE") == "PRODUCTION"
	addr := os.Getenv("ADDR")
	port := os.Getenv("PORT")

	staticPath := "./static"
	indexPath := "./index.html"
	if isProduction {
		staticPath = "/home/mwdev22/ins/static"
		indexPath = "/home/mwdev22/ins/index.html"
	}

	return Config{
		isProduction: isProduction,
		addr:         addr,
		staticPath:   staticPath,
		indexPath:    indexPath,
		Port:         port,
	}
}

func initializeRouter(config *Config) *http.ServeMux {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(config.staticPath))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		http.ServeFile(w, r, config.indexPath)
	})
	mux.HandleFunc("/calculate", calculate)
	mux.HandleFunc("/test", algTest)

	return mux
}

func startServer(config *Config, mux *http.ServeMux) error {
	if config.addr != "" {
		addr := config.addr + ":" + config.Port
		fmt.Println(addr)
		return http.ListenAndServe(addr, mux)
	}
	return http.ListenAndServe(":8080", restrictPaths(mux.ServeHTTP))
}

func modeName(isProduction bool) string {
	if isProduction {
		return "production"
	}
	return "development"
}

// restrykcja ścieżek
func restrictPaths(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allowedPaths := []string{"/", "/calculate", "/static/", "/selection", "/mutation", "/crossover"}
		for _, path := range allowedPaths {
			if r.URL.Path == path || (path == "/static/" && r.URL.Path[:8] == "/static/") {
				next.ServeHTTP(w, r)
				return
			}
		}
		http.NotFound(w, r)
	}
}
