package main

import "net/http"

type CalculationPayload struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
	D float64 `json:"d"`
	N int     `json:"N"`
}

type SelectionPayload struct {
	Population []Individual `json:"pop"`
	GSum       float64      `json:"g_sum"`
	A          float64      `json:"a"`
	B          float64      `json:"b"`
}

type CrossoverPayload struct {
	Population []Individual `json:"pop"`
	Pk         float64      `json:"pk"`
}

type Result struct {
	L          int          `json:"L"`
	Population []Individual `json:"population"`
	GSum       float64      `json:"g_sum"`
}

type SelectionResult struct {
	Population []Individual `json:"population"`
}

type Individual struct {
	ID      int     `json:"id"`
	XReal   float64 `json:"x_real"`
	XInt    int     `json:"x_int,omitempty"`
	Bin     string  `json:"bin"`
	Fx      float64 `json:"fx"`
	Gx      float64 `json:"gx"`
	P       float64 `json:"p,omitempty"`
	Q       float64 `json:"q"`
	R       float64 `json:"r"`
	XSel    float64 `json:"x_sel"`
	XSelBin string  `json:"x_sel_bin"`
	Parent  any     `json:"parent"`
}

type MutationPayload struct {
	Offspring    []Individual `json:"offspring"`
	MutationRate float64      `json:"mutation_rate"`
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
