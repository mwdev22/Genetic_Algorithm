package main

import "net/http"

// struktura potrzebna do zdekodowania danych od użytkownika
type CalculationPayload struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
	D float64 `json:"d"`
	N int     `json:"N"`
}

// struktura reprezentująca osobnika
type Individual struct {
	ID       int     `json:"id"`
	XReal    float64 `json:"x_real"`
	XInt     int     `json:"x_int"`
	Bin      string  `json:"bin"`
	XNewInt  int     `json:"x_new_int"`
	XNewReal float64 `json:"x_new_real"`
	Fx       float64 `json:"fx"`
}

// lista osobników oraz ich parametrów
type Result struct {
	Population []Individual `json:"population"`
	BestInd    Individual   `json:"best_ind"` // najlepiej dopasowany osobnik
	L          int          `json:"L"`
}

// restrykcja ścieżek
func restrictPaths(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allowedPaths := []string{"/", "/calculate", "/static/"}
		for _, path := range allowedPaths {
			if r.URL.Path == path || (path == "/static/" && r.URL.Path[:8] == "/static/") {
				next.ServeHTTP(w, r)
				return
			}
		}
		http.NotFound(w, r)
	}
}
