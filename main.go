package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var l int              // liczba bitów
var production = false // true w przy udostępnianiu na serwerze

func main() {
	fmt.Println("Starting app...")
	if os.Getenv("MODE") == "PRODUCTION" {
		production = true
	}
	rand.Seed(time.Now().UnixNano())

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
	mux.HandleFunc("/crossover", crossover)

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

	var result CalculationResult
	result.L = l
	for i := 1; i <= N; i++ {
		xReal := math.Round((a+rand.Float64()*(b-a))/d) * d
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

	var individuals []Individual = payload.Population

	var pSum float64 = 0

	for i := 0; i < len(individuals); i++ {
		indiv := &individuals[i]
		indiv.P = indiv.Gx / payload.GSum
		pSum += indiv.P
		indiv.Q = pSum
	}

	for i := 0; i < len(individuals); i++ {
		indiv := &individuals[i]
		indiv.R = rand.Float64()
		for j := 0; j < len(individuals); j++ {
			var qLast float64
			if j == 0 {
				qLast = 0
			} else {
				qLast = individuals[j-1].Q
			}
			if indiv.R > qLast && indiv.R < individuals[j].Q {
				indiv.XSel = individuals[j].XReal
				indiv.XSelBin = intToBin(realToInt(indiv.XSel, payload.A, payload.B))
			}
		}
	}

	result := SelectionResult{
		Population: individuals,
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, fmt.Sprintf("błąd przy encodingu odpowiedzi: %s", err), http.StatusInternalServerError)
	}

}

func crossover(w http.ResponseWriter, r *http.Request) {

	var payload CrossoverPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("błąd przy parsowaniu body requestu: %s", err), http.StatusBadRequest)
		return
	}

	var individuals []Individual = payload.Population

	for i := 0; i < len(individuals); i++ {
		ind := &individuals[i]
		r := rand.Float64()
		if r <= payload.Pk {
			ind.Parent = ind.XSelBin
		} else {
			ind.Parent = "-"
			ind.Child = "-"
			ind.NewGen = ind.XSelBin
			continue
		}
		ind.Pc = rand.Intn(len(ind.XSelBin)-2) + 1
	}

	var backupInd *Individual
	var backupId, backupPc int
	for i := 0; i < len(individuals); i++ {
		ind := &individuals[i]
		if ind.Parent == "-" || ind.Child != "" {
			continue
		}

		var secondParent *Individual
		for j := i + 1; j < len(individuals); j++ {
			secInd := &individuals[j]
			if secInd.Parent != "-" {
				secondParent = secInd
				break
			}
		}

		if secondParent == nil {
			ind.Child = ind.Parent[:ind.Pc] + backupInd.XSelBin[ind.Pc:]
			ind.NewGen = ind.Child
			backupId = backupInd.ID
			backupPc = ind.Pc
			continue
		}

		ind.Child = ind.Parent[:ind.Pc] + secondParent.XSelBin[ind.Pc:]
		ind.NewGen = ind.Child
		secondParent.Pc = ind.Pc
		secondParent.Child = secondParent.Parent[:secondParent.Pc] + ind.Parent[secondParent.Pc:]
		secondParent.NewGen = secondParent.Child

		backupInd = secondParent
	}
	result := CrossoverResult{
		Population: individuals,
		BackupId:   backupId,
		BackupPc:   backupPc,
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, fmt.Sprintf("błąd przy encodingu odpowiedzi: %s", err), http.StatusInternalServerError)
	}
}

func mutation(w http.ResponseWriter, r *http.Request) {

}

func g(xReal, a, b, d float64) float64 {
	return evalFunc(xReal) - minF(a, b, d) + d
}

func minF(a, b, d float64) float64 {
	var min float64 = evalFunc(a)
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
