package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting app...")

	// router
	mux := http.NewServeMux()

	mux.HandleFunc("/", index)
	mux.HandleFunc("/calculate", calculate)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("failed to start the server")
	}
}

func index(w http.ResponseWriter, r *http.Request) {

}

func calculate(w http.ResponseWriter, r *http.Request) {

}

func realToInt() {

}

func intToBin() {

}

func binToInt() {

}

func intToReal() {

}
