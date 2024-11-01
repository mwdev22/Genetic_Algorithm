package main

import "net/http"

// payloads
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

type MutationPayload struct {
	Population []Individual `json:"pop"`
	Pm         float64      `json:"pm"`
	A          float64      `json:"a"`
	B          float64      `json:"b"`
	D          float64      `json:"d"`
}

// responses
type CalculationResult struct {
	L          int          `json:"L"`
	Population []Individual `json:"population"`
	GSum       float64      `json:"g_sum"`
}

type PopulationResult struct {
	Population []Individual `json:"population"`
}

type CrossoverResult struct {
	Population []Individual `json:"population"`
	BackupId   int          `json:"backup_id"`
	BackupPc   int          `json:"backup_pc"`
}

type Individual struct {
	ID           int     `json:"id"`
	XReal        float64 `json:"x_real"`
	XInt         int     `json:"x_int,omitempty"`
	Bin          string  `json:"bin"`
	Fx           float64 `json:"fx"`
	Gx           float64 `json:"gx"`
	P            float64 `json:"p,omitempty"`
	Q            float64 `json:"q"`
	R            float64 `json:"r"`
	XSel         float64 `json:"x_sel"`
	XSelBin      string  `json:"x_sel_bin"`
	Parent       string  `json:"parent"`
	Pc           int     `json:"pc"`
	Child        string  `json:"child"`
	NewGen       string  `json:"new_gen"`
	MutatedGenes []int   `json:"mutated_genes"`
	FinalGen     string  `json:"final_gen"`
	FinalXReal   float64 `json:"final_x_real"`
	FinalFx      float64 `json:"final_fx"`
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
