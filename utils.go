package main

import (
	"net/http"
	"os"
)

type Config struct {
	isProduction bool
	addr         string
	staticPath   string
	indexPath    string
}

func loadConfig() Config {
	isProduction := os.Getenv("MODE") == "PRODUCTION"
	addr := os.Getenv("IP")

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
	}
}

func initializeRouter(config *Config) *http.ServeMux {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(config.staticPath))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, config.indexPath)
	})
	mux.HandleFunc("/calculate", calculate)
	mux.HandleFunc("/selection", selection)
	mux.HandleFunc("/crossover", crossover)
	mux.HandleFunc("/mutation", mutation)

	return mux
}

func startServer(config *Config, mux *http.ServeMux) error {
	if config.addr != "" {
		return http.ListenAndServe(config.addr, restrictPaths(mux.ServeHTTP))
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
