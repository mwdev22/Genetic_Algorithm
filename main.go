package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
)

var l int // liczba bitów

func main() {
	fmt.Println("Starting app...")

	// router
	mux := http.NewServeMux()

	mux.HandleFunc("/", index)
	mux.HandleFunc("/calculate", calculate)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("nie udało się uruchomić serwera: %v", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// funkcja dopasowania, trafności osobnika
func fitness(x float64) float64 {
	mod := x - math.Floor(x)
	return mod * (math.Cos(20*math.Pi*x) - math.Sin(x))
}

func calculate(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var payload CalculationPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "nieodpowiednie body requestu", http.StatusBadRequest)
		return
	}

	a := payload.A
	b := payload.B
	d := payload.D
	N := payload.N

	// kalkulacja liczby bitów
	l = int(math.Ceil(math.Log2((b - a) / d)))

	var result Result
	for i := 0; i < N; i++ {
		xReal := a + rand.Float64()*(b-a)
		xInt := realToInt(xReal, a, b)
		bin := intToBin(xInt)
		fx := fitness(xReal)

		indiv := Individual{
			XReal: xReal,
			XInt:  xInt,
			Bin:   bin,
			Fx:    fx,
		}
		result.Population = append(result.Population, indiv)
	}

	// Return the result as JSON
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
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
	return fmt.Sprintf("%b", xInt)
}
