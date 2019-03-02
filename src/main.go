package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	utils "./utils"

	"github.com/smartystreets/scanners/csv"
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

	readCsvFile("../csv/fullTest.csv")
	fmt.Println("END")

}

func readCsvFile(filePath string) {
	poolsize := 20
	jobch := make(chan utils.PhoneNumber)
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
		scanner, err := csv.NewStructScanner(bufio.NewReader(f))
		if err != nil {
			log.Panic(err)
		}
		var pn utils.PhoneNumber
		for scanner.Scan() {
			if err := scanner.Populate(&pn); err != nil {
				log.Panic(err)
			}
			jobch <- pn
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

func processData(jobs <-chan utils.PhoneNumber, results chan<- utils.PhoneNumber, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		validateStruct(&j)
		results <- j
	}
}
