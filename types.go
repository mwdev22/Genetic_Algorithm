package main

// payloads
type CalculationPayload struct {
	A     float64 `json:"a"`
	B     float64 `json:"b"`
	D     float64 `json:"d"`
	N     int     `json:"N"`
	T     int     `json:"T"`
	Pk    float64 `json:"pk"`
	Pm    float64 `json:"pm"`
	Elite bool    `json:"elite"`
}

type TestPayload struct {
	A     float64 `json:"a"`
	B     float64 `json:"b"`
	D     float64 `json:"d"`
	N     int     `json:"N"`
	T     int     `json:"T"`
	Pk    float64 `json:"pk"`
	Pm    float64 `json:"pm"`
	Elite bool    `json:"elite"`
}

// responses
type CalculationResult struct {
	Population   []*Individual      `json:"population"`
	GenStats     []*GenerationStats `json:"gen_stats"`
	FinalGenData []*FinalResult     `json:"results"`
}

type GenerationStats struct {
	FMin     float64 `json:"f_min"`
	FMax     float64 `json:"f_max"`
	FAvg     float64 `json:"f_avg"`
	Elite    float64
	EliteInd int
}

type FinalResult struct {
	XReal   float64 `json:"x_real"`
	XBin    string  `json:"x_bin"`
	Fx      float64 `json:"fx"`
	Percent float64 `json:"percent"`
	Count   int     `json:"count"`
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
