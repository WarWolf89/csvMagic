package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	root "../../pkg"
	mongoutils "../../pkg/mongoutils"
	validutils "../../pkg/validutils"

	"github.com/smartystreets/scanners/csv"
)

var (
	jobch           = make(chan root.PhoneNumber)
	results         = make(chan root.PhoneNumber)
	client, context = mongoutils.SetupConnection("mongodb://localhost:27017")
)

func main() {

	readCsvFile("../../resources/csv/test.csv")
	fmt.Println("END")

}

func readCsvFile(filePath string) {
	csvService := mongoutils.CreateCsvService(client, "local", "csv-test")
	poolsize := 20
	var wg sync.WaitGroup
	start := time.Now()
	counter := 0
	// set up workers
	for w := 1; w <= poolsize; w++ {
		wg.Add(1)
		go processData(jobch, results, &wg, csvService)
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
		var pn root.PhoneNumber
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

func processData(jobs <-chan root.PhoneNumber, results chan<- root.PhoneNumber, wg *sync.WaitGroup, csvService *mongoutils.CsvService) {
	defer wg.Done()
	for j := range jobs {
		validutils.CheckAndFixStruct(&j)
		res, err := csvService.Collection.InsertOne(csvService.Context, j)
		fmt.Println(res)
		if err != nil {
			panic(err)
		}
		results <- j
	}
}
