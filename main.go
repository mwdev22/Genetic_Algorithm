package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var l int // liczba bitów

var testStarted = false
var resultsChan = make(chan *TestResult)
var done = make(chan struct{})

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

	var population []*Individual
	var statsSummary []*GenerationStats

	population, statsSummary = calculateGenerations(&payload)

	var results = make(map[float64]*FinalResult)

	for _, ind := range population {
		finalXReal := math.Round(ind.XReal/payload.D) * payload.D
		res, ok := results[finalXReal]
		if !ok {
			results[finalXReal] = &FinalResult{
				XReal:   finalXReal,
				XBin:    ind.FinalGen,
				Fx:      ind.FinalFx,
				Count:   1,
				Percent: (1.0 / float64(len(population))) * 100,
			}
		} else {
			res.Count += 1
			res.Percent = (float64(res.Count) / float64(len(population))) * 100
		}
	}

	var finalResults []*FinalResult
	for _, v := range results {
		finalResults = append(finalResults, v)
	}

	// formatowanie odpowiedzi do formatu JSON
	rsp := CalculationResult{
		Population:   population,
		GenStats:     statsSummary,
		FinalGenData: finalResults,
	}
	err = json.NewEncoder(w).Encode(rsp)
	if err != nil {
		http.Error(w, fmt.Sprintf("błąd przy encodingu odpowiedzi: %s", err), http.StatusInternalServerError)
	}
}

func calculateGenerations(payload *CalculationPayload) ([]*Individual, []*GenerationStats) {
	a := payload.A
	b := payload.B
	d := payload.D
	N := payload.N
	T := payload.T
	pk := payload.Pk
	pm := payload.Pm

	isElite := payload.Elite

	// Calculate the number of bits
	l = int(math.Ceil(math.Log2((b - a) / d)))
	minFx := minF(a, b, d)

	var individuals []*Individual
	var genStatsSummary []*GenerationStats

	var gSum float64 = 0

	for i := 1; i <= N; i++ {
		xReal := math.Round((a+rand.Float64()*(b-a))/d) * d
		xInt := realToInt(xReal, a, b)
		bin := intToBin(xInt)
		fx := evalFunc(xReal)
		gx := g(fx, minFx, d)

		gSum += gx

		indiv := &Individual{
			ID:    i,
			XReal: xReal,
			Bin:   bin,
			Fx:    fx,
			Gx:    gx,
		}
		individuals = append(individuals, indiv)
	}

	for i := 0; i < T; i++ {
		gSum = 0
		for _, ind := range individuals {
			gSum += ind.Gx
		}

		var genStats *GenerationStats
		individuals, genStats = genAlgorithm(a, b, d, pk, pm, minFx, gSum, individuals, isElite)
		genStatsSummary = append(genStatsSummary, genStats)
	}
	return individuals, genStatsSummary
}

func algTest(w http.ResponseWriter, r *http.Request) {
	if !testStarted {

		a, b, d := -4.0, 12.0, 0.001
		N := []int{30, 35, 40, 45, 50, 55, 60, 65, 70, 75, 80}
		pk := []float64{0.5, 0.55, 0.6, 0.65, 0.7, 0.75, 0.8, 0.85, 0.9}
		T := []int{50, 60, 70, 80, 90, 100}
		pm := []float64{0.0001, 0.0005, 0.001, 0.005, 0.01}

		testStarted = true

		go runTest(a, b, d, N, pk, T, pm)
		w.Write([]byte("test started"))
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	go func() {
		var bestResults = make(map[float64]TestResult)

		for result := range resultsChan {
			bestRes, exists := bestResults[result.FAvg]
			sendResult := true

			if exists {

				currentParams := result.T + result.N
				bestParams := bestRes.T + bestRes.N
				if currentParams >= bestParams {
					sendResult = false
				} else {
					bestResults[result.FAvg] = *result
				}
			}

			if sendResult {
				bestResults[result.FAvg] = *result
				data, err := json.Marshal(result)
				if err == nil {
					fmt.Fprintf(w, "data: %s\n\n", data)
					w.(http.Flusher).Flush()
				}
			}

			fmt.Println(result.FAvg, result.N, result.T)
		}

		done <- struct{}{}
	}()

	<-done
}

func runTest(a, b, d float64, N []int, pk []float64, T []int, pm []float64) {
	var wg sync.WaitGroup

	sem := make(chan struct{}, 16)

	for _, n := range N {
		for _, t := range T {
			for _, pkVal := range pk {
				for _, pmVal := range pm {
					wg.Add(1)

					sem <- struct{}{}

					go func(n, t int, pkVal, pmVal float64) {
						defer func() {
							<-sem
							wg.Done()
						}()

						var localBestResult TestResult
						for i := 0; i < 100; i++ {
							_, genStats := calculateGenerations(&CalculationPayload{
								A:     a,
								B:     b,
								D:     d,
								N:     n,
								Pk:    pkVal,
								Pm:    pmVal,
								T:     t,
								Elite: true,
							})

							if genStats[t-1].FAvg > localBestResult.FAvg {
								localBestResult = TestResult{
									N:    n,
									T:    t,
									Pk:   pkVal,
									Pm:   pmVal,
									FAvg: genStats[t-1].FAvg,
								}
							}
						}

						resultsChan <- &localBestResult
					}(n, t, pkVal, pmVal)
				}
			}
		}
	}

	wg.Wait()
	close(resultsChan)
}

func genAlgorithm(a, b, d, pk, pm, minFx, gSum float64, individuals []*Individual, isElite bool) ([]*Individual, *GenerationStats) {

	selection(gSum, a, b, individuals)
	crossover(individuals, pk)
	genStats := mutationAndStatsNote(pm, a, b, d, minFx, individuals, isElite)

	return individuals, genStats
}

func selection(gSum, a, b float64, individuals []*Individual) {

	var pSum float64 = 0

	for i := 0; i < len(individuals); i++ {
		indiv := individuals[i]
		indiv.P = indiv.Gx / gSum
		pSum += indiv.P
		indiv.Q = pSum
	}

	for i := 0; i < len(individuals); i++ {
		indiv := individuals[i]
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
				indiv.XSelBin = intToBin(realToInt(indiv.XSel, a, b))
			}
		}
	}

}

func crossover(individuals []*Individual, pk float64) {

	// incjalizacja danych krzyżowania osobnikóœ
	for i := 0; i < len(individuals); i++ {
		ind := individuals[i]
		r := rand.Float64()
		if r <= pk {
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
	for i := 0; i < len(individuals); i++ {
		ind := individuals[i]
		if ind.Parent == "-" || ind.Child != "" {
			continue
		}

		var secondParent *Individual
		for j := i + 1; j < len(individuals); j++ {
			secInd := individuals[j]
			if secInd.Parent != "-" {
				secondParent = secInd
				break
			}
		}

		if secondParent == nil {
			ind.Child = ind.Parent[:ind.Pc] + backupInd.XSelBin[ind.Pc:]
			ind.NewGen = ind.Child
			continue
		}

		ind.Child = ind.Parent[:ind.Pc] + secondParent.XSelBin[ind.Pc:]
		ind.NewGen = ind.Child
		secondParent.Pc = ind.Pc
		secondParent.Child = secondParent.Parent[:secondParent.Pc] + ind.Parent[secondParent.Pc:]
		secondParent.NewGen = secondParent.Child

		backupInd = secondParent
	}
}

func mutationAndStatsNote(pm, a, b, d, minFx float64, individuals []*Individual, isElite bool) *GenerationStats {

	var fMin, fMax, fSum, fAvg, bestXReal, bestFx float64

	for i := 0; i < len(individuals); i++ {
		ind := individuals[i]
		finalGen := []byte(ind.NewGen)
		for j := 0; j < len(finalGen); j++ {
			r := rand.Float64()
			if r <= pm {
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
		ind.FinalXReal = intToReal(binToInt(ind.FinalGen), a, b)
		ind.FinalXReal = math.Round(ind.FinalXReal/d) * d
		ind.FinalFx = evalFunc(ind.FinalXReal)

		if ind.Fx > bestFx {
			bestFx = ind.Fx
			bestXReal = ind.XReal
		}
	}

	randVal := rand.Intn(len(individuals))
	randomInd := individuals[randVal]
	if isElite && randomInd.FinalFx < bestFx {
		randomInd.XReal = bestXReal
		randomInd.Fx = bestFx
		randomInd.FinalFx = bestFx
		randomInd.FinalXReal = bestXReal
		randomInd.XInt = realToInt(bestXReal, a, b)
		randomInd.Bin = intToBin(randomInd.XInt)
		randomInd.Gx = g(randomInd.Fx, minFx, d)
	}

	for _, ind := range individuals {
		ind.XReal = ind.FinalXReal
		ind.XInt = realToInt(ind.XReal, a, b)
		ind.Bin = ind.FinalGen
		ind.Fx = ind.FinalFx
		ind.Gx = g(ind.Fx, minFx, d)

		if ind.FinalFx < fMin {
			fMin = ind.FinalFx
		}
		if ind.FinalFx > fMax {
			fMax = ind.FinalFx
		}
		fSum += ind.FinalFx
	}

	fAvg = fSum / float64(len(individuals))

	genStats := &GenerationStats{
		FMin:     fMin,
		FMax:     fMax,
		FAvg:     fAvg,
		Elite:    bestFx,
		EliteInd: randVal,
	}

	return genStats

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
