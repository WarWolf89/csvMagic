package main

import (
	"fmt"
	"log"
	"net/http"

	handlers "../../pkg/handlers"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/csvUpload", handlers.CsvUpload).Methods("POST")
	router.HandleFunc("/valSingle", handlers.ValidateSingleNumber).Methods("POST")
	fmt.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
	fmt.Println("END")
}
