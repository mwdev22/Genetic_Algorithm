package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Starting app...")

	http.HandleFunc("/", index)
	http.HandleFunc("/calculate", calculate)
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
