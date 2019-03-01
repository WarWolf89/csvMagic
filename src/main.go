package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	utils "./utils"

	validator "gopkg.in/go-playground/validator.v9"
)

func validateStruct(pn *utils.PhoneNumber) {
	validate.RegisterValidation("custom", utils.ValidateFieldForSMSPhone)

	err := validate.Struct(pn)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}

		utils.FixVal(pn, err)
	}
}

var validate *validator.Validate

func main() {

	validate = validator.New()

	readCsvFile("../csv/test.csv")

	fmt.Println("END")

}

func readCsvFile(filePath string) {
	poolsize := 20
	jobch := make(chan []string)
	results := make(chan utils.PhoneNumber)
	var wg sync.WaitGroup
	start := time.Now()
	counter := 0
	// set up workers
	for w := 1; w <= poolsize; w++ {
		wg.Add(1)
		go processData(jobch, results, &wg)
	}
	// Load a csv file.
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	// Create a new reader.
	go func() {
		reader := csv.NewReader(bufio.NewReader(f))
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			jobch <- record
		}
		close(jobch)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for v := range results {
		fmt.Println(v)
		counter++
	}
	fmt.Println(counter)
	fmt.Printf("\n%2fs", time.Since(start).Seconds())
}

func processData(jobs <-chan []string, results chan<- utils.PhoneNumber, wg *sync.WaitGroup) {
	defer wg.Done()
	var pn utils.PhoneNumber
	for j := range jobs {
		pn = utils.PhoneNumber{j[0], j[1]}
		validateStruct(&pn)
		results <- pn
	}

}
