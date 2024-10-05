package main

type CalculationPayload struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
	D float64 `json:"d"`
	N int     `json:"N"`
}

type Result struct {
	Population []Individual `json:"population"`
}

type Individual struct {
	XReal float64 `json:"x_real"`
	XInt  int     `json:"x_int"`
	Bin   string  `json:"bin"`
	Fx    float64 `json:"fx"`
}
