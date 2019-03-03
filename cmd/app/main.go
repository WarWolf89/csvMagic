package main

import (
	"fmt"

	reader "../../pkg/reader"
)

func main() {
	reader.ReadCsvFile("../../resources/csv/test.csv")
	fmt.Println("END")
}
