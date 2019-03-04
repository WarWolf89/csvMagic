package reader

import (
	"bufio"
	"fmt"
	"log"
	"mime/multipart"
	"sync"

	root ".."
	mongoutils "../mongoutils"
	validutils "../validutils"
	"github.com/rs/xid"
	"github.com/smartystreets/scanners/csv"
)

var (
	client, context = mongoutils.SetupConnection("mongodb://localhost:27017")
)

func ReadCsvFile(file multipart.File, header string) (*root.FileMeta, error) {
	jobch := make(chan root.PhoneNumber)
	results := make(chan root.PhoneNumber)
	poolsize := 20
	var wg sync.WaitGroup

	// generate uuid for the file and use it as a reference for later lookups
	uuID := xid.New()
	// set up the meta struct for response
	fm := root.NewFileMeta(uuID.String(), header)
	// Set up the mongodb service
	metaService := mongoutils.CreateCsvService(client, "local", "META")
	csvService := mongoutils.CreateCsvService(client, "local", "csv-test")
	// set up workers for the pool
	for w := 1; w <= poolsize; w++ {
		wg.Add(1)
		go processData(jobch, results, &wg, fm)
	}
	// Create a new reader.
	go func() {
		scanner, err := csv.NewStructScanner(bufio.NewReader(file))
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

	writeRowsToMongo(results, csvService, fm)
	writeMetaToMongo(metaService, fm)
	return fm, nil
}

func processData(jobs <-chan root.PhoneNumber, results chan<- root.PhoneNumber, wg *sync.WaitGroup, fm *root.FileMeta) {
	defer wg.Done()

	for j := range jobs {
		validutils.CheckAndFixStruct(&j, fm)
		results <- j
	}
}

func writeMetaToMongo(metaService *mongoutils.CsvService, fm *root.FileMeta) {
	res, inErr := metaService.Collection.InsertOne(metaService.Context, fm)
	if inErr != nil {
		fmt.Println(inErr)
	}
}

func writeRowsToMongo(results <-chan root.PhoneNumber, csvService *mongoutils.CsvService, fm *root.FileMeta) {
	for r := range results {
		_, err := csvService.Collection.InsertOne(csvService.Context, r)
		if err != nil {
			fm.Errors = append(fm.Errors, err.Error())
			continue
		}
		fm.IncreaseCounter("processed")
	}
}
