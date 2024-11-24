package main

// payloads
type CalculationPayload struct {
	A    float64 `json:"a"`
	B    float64 `json:"b"`
	D    float64 `json:"d"`
	TMax int     `json:"T"`
}

type CalculationResponse struct {
	VcData     []*IterData `json:"vc_data"`
	MaxResults []*MaxStep  `json:"max_results"`
}
type TestResponse struct {
	StatsMap    map[int]int     `json:"stats_map"`
	Percentages map[int]float64 `json:"percentages"`
	SuccessSum  int             `json:"success_sum"`
}

type Vc struct {
	XReal float64 `json:"x_real"`
	XBin  string  `json:"x_bin"`
	Fx    float64 `json:"fx"`
}

type Vn struct {
	XReal float64 `json:"x_real"`
	XBin  string  `json:"x_bin"`
	Fx    float64 `json:"fx"`
}

type MaxStep struct {
	T     int
	MaxFx float64
}

type LocalStep struct {
	index int
	Fx    float64
}

type IterData struct {
	Vcs   []*Vc
	Steps []*LocalStep
}
