package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type phoneNumber struct {
	id       string
	smsPhone string
}

func main() {
	readCsvFile("../csv/South_African_Mobile_Numbers.csv")

}

func readCsvFile(filePath string) {

	start := time.Now()

	// Load a csv file.
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	// Create a new reader.
	reader := csv.NewReader(bufio.NewReader(f))
	var wg sync.WaitGroup
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		wg.Add(1)
		go func(record []string) {
			defer wg.Done()
			processData(record)
		}(record)
	}

	// closer
	wg.Wait()
	fmt.Printf("\n%2fs", time.Since(start).Seconds())
}

func processData(r []string) {

	fmt.Println("started job", r)
	fmt.Println("finished job", r)
}
