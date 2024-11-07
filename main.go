package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"
)

var l int // liczba bitów

func main() {
	config := loadConfig()
	log.Printf("uruchamianie aplikacji w trybie %s...", modeName(config.isProduction))

	rand.Seed(time.Now().UnixNano())

	mux := initializeRouter(&config)

	if err := startServer(&config, mux); err != nil {
		log.Fatalf("bład przy uruchamianiu serwera: %v", err)
	}
}

// funkcja oceny
func evalFunc(x float64) float64 {
	// mod := x - math.Floor(x)
	// return mod * (math.Cos(20*math.Pi*x) - math.Sin(x))
	return -(x + 1) * (x - 1) * (x - 2)
}

func calculate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var payload CalculationPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("błąd przy parsowaniu body requestu: %s", err), http.StatusBadRequest)
		return
	}

	a, b, d, N := payload.A, payload.B, payload.D, payload.N
	l = int(math.Ceil(math.Log2((b - a) / d)))
	minFx := minF(a, b, d)

	var result CalculationResult
	result.L = l
	gSumChan := make(chan float64, N)
	populationChan := make(chan Individual, N)

	for i := 1; i <= N; i++ {
		go func(id int) {
			xReal := math.Round((a+rand.Float64()*(b-a))/d) * d
			xInt := realToInt(xReal, a, b)
			bin := intToBin(xInt)
			fx := evalFunc(xReal)
			gx := g(fx, minFx, d)

			gSumChan <- gx
			populationChan <- Individual{
				ID:    id,
				XReal: xReal,
				Bin:   bin,
				Fx:    fx,
				Gx:    gx,
			}
		}(i)
	}

	var gSum float64
	for i := 0; i < N; i++ {
		gSum += <-gSumChan
		result.Population = append(result.Population, <-populationChan)
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

	result := PopulationResult{
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

	// incjalizacja danych krzyżowania osobnikóœ
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
		ind.Pc = rand.Intn(len(ind.XSelBin)-1) + 1
	}

	var backupInd *Individual
	// przechowywanie zapasowego rodzica, w razie nieparzystej ilości
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

	var payload MutationPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("błąd przy parsowaniu body requestu: %s", err), http.StatusBadRequest)
		return
	}

	var individuals []Individual = payload.Population

	for i := 0; i < len(individuals); i++ {
		ind := &individuals[i]
		finalGen := []byte(ind.NewGen)
		for j := 0; j < len(finalGen); j++ {
			r := rand.Float64()
			if r <= payload.Pm {
				switch finalGen[j] {
				case '0':
					finalGen[j] = '1'
				case '1':
					finalGen[j] = '0'
				}
				ind.MutatedGenes = append(ind.MutatedGenes, j+1)
			}
		}
		ind.FinalGen = string(finalGen)
		ind.FinalXReal = intToReal(binToInt(ind.FinalGen), payload.A, payload.B)
		ind.FinalXReal = math.Round(ind.FinalXReal/payload.D) * payload.D
		ind.FinalFx = evalFunc(ind.FinalXReal)
	}

	result := PopulationResult{
		Population: individuals,
	}
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		http.Error(w, fmt.Sprintf("błąd przy encodingu odpowiedzi: %s", err), http.StatusInternalServerError)
	}

}

func g(fx, minFx, d float64) float64 {
	return fx - minFx + d
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
