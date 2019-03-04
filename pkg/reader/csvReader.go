package reader

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	root ".."
	mongoutils "../mongoutils"
	validutils "../validutils"
	"github.com/rs/xid"
	"github.com/smartystreets/scanners/csv"
)

var (
	jobch           = make(chan root.PhoneNumber)
	results         = make(chan root.PhoneNumber)
	client, context = mongoutils.SetupConnection("mongodb://localhost:27017")
)

func ReadCsvFile(filePath string) (*root.FileMeta, error) {
	poolsize := 20
	var wg sync.WaitGroup
	start := time.Now()
	// Load a csv file.
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("no such file")
	}
	defer f.Close()
	// generate uuid for the file and use it as a reference for later lookups
	uuID := xid.New()
	// set up the meta struct for response
	fm := root.NewFileMeta(uuID.String(), f.Name())
	// Set up the mongodb service
	metaService := mongoutils.CreateCsvService(client, "local", "META")

	colln := "abc"
	csvService := mongoutils.CreateCsvService(client, "local", colln)

	// set up workers for the pool
	for w := 1; w <= poolsize; w++ {
		wg.Add(1)
		go processData(jobch, results, &wg, fm)
	}
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

			pn.FileID = uuID.String()
			jobch <- pn
		}
		close(jobch)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	writeToMongo(results, csvService, fm)
	fm.ExecTime = time.Since(start).Seconds()

	_, inErr := metaService.Collection.InsertOne(metaService.Context, fm)
	if inErr != nil {
		fmt.Println(inErr)
	}
	return fm, nil
}

func processData(jobs <-chan root.PhoneNumber, results chan<- root.PhoneNumber, wg *sync.WaitGroup, fm *root.FileMeta) {
	defer wg.Done()
	for j := range jobs {
		validutils.CheckAndFixStruct(&j, fm)
		results <- j
	}
}

func writeToMongo(results <-chan root.PhoneNumber, csvService *mongoutils.CsvService, fm *root.FileMeta) {
	for r := range results {
		_, err := csvService.Collection.InsertOne(csvService.Context, r)
		if err != nil {
			fm.Errors = append(fm.Errors, err.Error())
			continue
		}
		fm.IncreaseCounter("processed")
	}
}
