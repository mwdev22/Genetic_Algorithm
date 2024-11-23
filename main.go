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
func F(x float64) float64 {
	mod := x - math.Floor(x)
	return mod * (math.Cos(20*math.Pi*x) - math.Sin(x))
	// return -(x + 1) * (x - 1) * (x - 2)
}

func calculate(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var payload CalculationPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("błąd przy parsowaniu body requestu: %s", err), http.StatusBadRequest)
		return
	}

	l = int(math.Ceil(math.Log2((payload.B - payload.A) / payload.D)))

	vcData, maxResults := findMaximumGrowth(&payload)
	var rsp = CalculationResponse{
		VcData:     vcData,
		MaxResults: maxResults,
	}
	err = json.NewEncoder(w).Encode(rsp)
	if err != nil {
		http.Error(w, fmt.Sprintf("błąd przy encodingu odpowiedzi: %s", err), http.StatusInternalServerError)
	}
}
func findMaximumGrowth(data *CalculationPayload) ([]*IterData, []*MaxStep) {
	var local bool = false
	var itersData []*IterData

	var maxResults []*MaxStep
	var maxVal float64

	for i := 0; i < data.TMax; i++ {
		local = false
		randomNumber := rand.Float64()*(data.B-data.A) + data.A
		// Round the random number to the nearest multiple of payload.D
		randomNumber = roundToNearest(randomNumber, data.D)

		xInt := realToInt(randomNumber, data.A, data.B)
		xBin := intToBin(xInt)
		vc := &Vc{
			XReal: randomNumber,
			XBin:  xBin,
			Fx:    F(randomNumber),
		}
		iterVcs := []*Vc{vc}
		firstStep := &LocalStep{index: 1, Fx: vc.Fx}
		localSteps := []*LocalStep{firstStep}
		if i == 0 {
			maxResults = append(maxResults, &MaxStep{
				MaxFx: vc.Fx,
				T:     i,
			})
		}
		for !local {
			bestNeigh := generateBestNeighbor(data.A, data.B, data.D, vc.XBin)
			fmt.Println(bestNeigh.Fx)
			if bestNeigh.Fx > vc.Fx {
				vc = &Vc{
					XBin:  bestNeigh.XBin,
					XReal: bestNeigh.XReal,
					Fx:    bestNeigh.Fx,
				}
				iterVcs = append(iterVcs, vc)
				newStep := &LocalStep{
					index: len(localSteps),
					Fx:    vc.Fx,
				}
				localSteps = append(localSteps, newStep)
			} else {
				fmt.Println(vc.Fx, bestNeigh.Fx)
				if vc.Fx > maxVal {
					maxVal = vc.Fx
				}
				local = true
			}
		}
		iterationData := &IterData{
			Vcs:   iterVcs,
			Steps: localSteps,
		}
		maxIterStep := &MaxStep{
			T:     i + 1,
			MaxFx: maxVal,
		}
		maxResults = append(maxResults, maxIterStep)
		itersData = append(itersData, iterationData)
	}

	return itersData, maxResults
}

func generateBestNeighbor(a, b, d float64, vc string) *Vn {
	var bestNeigh Vn
	for i := 0; i < len(vc); i++ {
		newBin := []rune(vc)
		switch vc[i] {
		case '1':
			newBin[i] = '0'
		case '0':
			newBin[i] = '1'
		}
		xReal := intToReal(binToInt(string(newBin)), a, b)
		// Round the XReal value to the nearest multiple of d
		xReal = roundToNearest(xReal, d)
		newV := &Vn{
			XReal: xReal,
			XBin:  string(newBin),
			Fx:    F(xReal),
		}

		if newV.Fx > bestNeigh.Fx {
			bestNeigh = *newV
		}
	}
	return &bestNeigh
}

func roundToNearest(value, step float64) float64 {
	return math.Round(value/step) * step
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
