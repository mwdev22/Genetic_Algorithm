package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
)

var l int              // liczba bitów
var production = false // true w przy udostępnianiu na serwerze

func main() {
	fmt.Println("Starting app...")
	if os.Getenv("MODE") == "PRODUCTION" {
		production = true
	}

	// router
	mux := http.NewServeMux()

	var staticPath string
	if production {
		staticPath = "/home/mwdev22/ins/static"
	} else {
		staticPath = "./static"
	}
	fs := http.FileServer(http.Dir(staticPath))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", index)
	mux.HandleFunc("/calculate", calculate)
	mux.HandleFunc("/selection", selection)

	var err error
	addr := os.Getenv("ADDR")
	if addr != "" {
		// wymaganie serwisu hostingowego do nasłuchiwania na ipv6
		err = http.ListenAndServe(addr, mux)
	} else {
		err = http.ListenAndServe(":8080", restrictPaths(mux.ServeHTTP))
	}
	if err != nil {
		log.Fatalf("nie udało się uruchomić serwera: %v", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	var path string
	if production {
		path = "/home/mwdev22/ins/index.html"
	} else {
		path = "./index.html"
	}
	http.ServeFile(w, r, path)
}

// funkcja oceny
func evalFunc(x float64) float64 {
	mod := x - math.Floor(x)
	return mod * (math.Cos(20*math.Pi*x) - math.Sin(x))
}

func calculate(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var payload CalculationPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("błąd przy parsowaniu body requestu: %s", err), http.StatusBadRequest)
		return
	}

	a := payload.A
	b := payload.B
	d := payload.D
	N := payload.N

	// kalkulacja liczby bitów
	l = int(math.Ceil(math.Log2((b - a) / d)))

	var gSum float64 = 0

	var result Result
	result.L = l
	for i := 1; i <= N; i++ {
		xReal := a + rand.Float64()*(b-a)
		xInt := realToInt(xReal, a, b)
		bin := intToBin(xInt)
		fx := evalFunc(xReal)
		gx := g(xReal, a, b, d)

		gSum += gx

		indiv := Individual{
			ID:    i,
			XReal: xReal,
			Bin:   bin,
			Fx:    fx,
			Gx:    gx,
		}
		result.Population = append(result.Population, indiv)

	}
	result.GSum = gSum

	// formatowanie odpowiedzi do formatu JSON
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, fmt.Sprintf("błąd przy encodingu odpowiedzi: %s", err), http.StatusInternalServerError)
	}
}

func selection(w http.ResponseWriter, r *http.Request) {
	var payload SelectionPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("błąd przy parsowaniu body requestu: %s", err), http.StatusBadRequest)
		return
	}
	var indinviduals = payload.Population

	for indiv := range indinviduals {
		// q :=
	}

	err = json.NewEncoder(w).Encode(payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("błąd przy encodingu odpowiedzi: %s", err), http.StatusInternalServerError)
	}

}

func mutation(x string, genIdxs []int) {

}

func crossover(x, y string, pc int) {

}

func g(xReal, a, b, d float64) float64 {
	return evalFunc(xReal) - minF(a, b, d) + d
}

func minF(a, b, d float64) float64 {
	min := evalFunc(a)
	x := a
	for x <= b {
		fx := evalFunc(x)
		if fx < min {
			min = fx
		}
		x += d
	}
	return min
}

// konwersje liczb
func realToInt(x, a, b float64) int {
	return int((x - a) / (b - a) * (math.Pow(2, float64(l)) - 1))
}

func intToReal(xInt int, a, b float64) float64 {
	return float64(xInt)*(b-a)/(math.Pow(2, float64(l))-1) + a
}

func intToBin(xInt int) string {
	return fmt.Sprintf("%0*b", l, xInt)
}

func binToInt(bin string) int {
	var xInt int
	fmt.Sscanf(bin, "%b", &xInt)
	return xInt
}
