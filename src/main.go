package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"sync"
	"time"

	validator "gopkg.in/go-playground/validator.v9"
)

type phoneNumber struct {
	id       string
	SmsPhone string `validate:"min=11,max=11"`
}

func validateStruct(pn *phoneNumber) {

	err := validate.Struct(pn)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}

		for _, err := range err.(validator.ValidationErrors) {
			val := fmt.Sprintf("%v", err.Value())
			field := err.Field()
			fmt.Printf("value before fixed %s \n ", val)
			fixVal(&val)
			reflect.ValueOf(pn).Elem().FieldByName(field).SetString(val)
			fmt.Printf("value after fixed %s \n ", pn.SmsPhone)
		}

		return
	}
}

func fixVal(valToFix *string) {
	re := regexp.MustCompile("\\D")
	*valToFix = re.ReplaceAllString(*valToFix, "")
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
	results := make(chan phoneNumber)
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

func processData(jobs <-chan []string, results chan<- phoneNumber, wg *sync.WaitGroup) {
	defer wg.Done()
	var pn phoneNumber
	for j := range jobs {
		pn = phoneNumber{j[0], j[1]}
		validateStruct(&pn)
		results <- pn
	}

}
