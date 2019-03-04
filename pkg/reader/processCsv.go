package reader

import (
	"bufio"
	"fmt"
	"log"
	"mime/multipart"
	"sync"
	"time"

	root ".."
	mongoutils "../mongoutils"
	validutils "../validutils"
	"github.com/rs/xid"
	"github.com/smartystreets/scanners/csv"
)

var (
	client, context = mongoutils.SetupConnection("mongodb://localhost:27017")
)

func ProcessCsv(file multipart.File, name string) (*root.FileMeta, error) {
	jobch := make(chan root.PhoneNumber)
	poolsize := 20
	start := time.Now()
	var wg sync.WaitGroup

	// generate uuid for the file and use it as a reference for later lookups
	uuID := xid.New()
	// set up the meta struct for response
	fm := root.NewFileMeta(uuID.String(), name)
	// Set up the mongodb service
	metaService := mongoutils.CreateCsvService(client, "local", "META")
	csvService := mongoutils.CreateCsvService(client, "local", "csv-test")
	// set up workers for the pool
	for w := 1; w <= poolsize; w++ {
		wg.Add(1)
		go processData(jobch, &wg, csvService, fm)
	}
	// Create a new reader.

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

	// wait for all jobs to finish
	wg.Wait()
	// set execution time and write file data to mongo
	fm.ExecTime = time.Since(start).Seconds()
	writeMetaToMongo(metaService, fm)

	return fm, nil
}

func processData(jobs <-chan root.PhoneNumber, wg *sync.WaitGroup, csvService *mongoutils.CsvService, fm *root.FileMeta) {
	defer wg.Done()
	for j := range jobs {
		validutils.CheckAndFixStruct(&j, fm)
		_, err := csvService.Collection.InsertOne(csvService.Context, j)
		if err != nil {
			fm.Errors = append(fm.Errors, err.Error())
			continue
		}
		fm.IncreaseCounter("processed")
	}
}

func writeMetaToMongo(metaService *mongoutils.CsvService, fm *root.FileMeta) *root.FileMeta {
	_, inErr := metaService.Collection.InsertOne(metaService.Context, fm)
	if inErr != nil {
		fmt.Println(inErr)
	}
	return fm
}
