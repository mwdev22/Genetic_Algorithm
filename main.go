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

	var result Result
	result.L = l
	for i := 1; i <= N; i++ {
		xReal := a + rand.Float64()*(b-a)
		xInt := realToInt(xReal, a, b)
		bin := intToBin(xInt)
		xNewInt := binToInt(bin)
		xNewReal := intToReal(xNewInt, a, b)

		fx := evalFunc(xNewReal)

		indiv := Individual{
			ID:       i,
			XReal:    xReal,
			XInt:     xInt,
			Bin:      bin,
			XNewInt:  xNewInt,
			XNewReal: xNewReal,
			Fx:       fx,
		}
		result.Population = append(result.Population, indiv)
		if fx > result.BestInd.Fx {
			result.BestInd = indiv
		}

	}

	// formatowanie odpowiedzi do formatu JSON
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, fmt.Sprintf("błąd przy encodingu odpowiedzi: %s", err), http.StatusInternalServerError)
	}
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
