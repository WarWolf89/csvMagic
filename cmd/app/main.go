package main

import (
	"fmt"
	"log"
	"net/http"

	"../../pkg/handlers"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/csvUpload", handlers.CsvUpload).Methods("POST")
	router.HandleFunc("/valSingle/{num}", handlers.ValidateSingleNumber).Methods("GET")
	router.HandleFunc("/filedata/{fid}", handlers.GetFileInfo).Methods("GET")
	fmt.Println("Starting server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))

}
