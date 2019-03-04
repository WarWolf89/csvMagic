package main

import (
	"fmt"

	reader "../../pkg/reader"
)

func main() {
	reader.ReadCsvFile("../../resources/csv/fullTest.csv")
	fmt.Println("END")
}
