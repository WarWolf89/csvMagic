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
	"github.com/smartystreets/scanners/csv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	addr    = "mongodb://localhost:27017"
	db      = "local"
	iCSV    = "file_id"
	collCSV = "csv-test"
	collMet = "META"
)

var (
	metaService = mongoutils.NewDBService(addr, db, collMet)
	csvService  = mongoutils.NewDBService(addr, db, collCSV)
)

func init() {
	iName, err := csvService.PopulateIndex(iCSV)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("created index %s \n", *iName)
}

func ProcessCsv(file multipart.File, name string) (*root.FileMeta, error) {
	jobch := make(chan root.PhoneNumber)
	poolsize := 20
	start := time.Now()
	var wg sync.WaitGroup
	var pn root.PhoneNumber
	// generate uuid for the file and use it as a reference for later lookups
	ID := primitive.NewObjectID()
	// set up the meta struct for response
	fm := root.NewFileMeta(ID, name)
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
	for scanner.Scan() {
		if err := scanner.Populate(&pn); err != nil {
			log.Panic(err)
		}
		pn.FileID = ID.Hex()
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

func processData(jobs <-chan root.PhoneNumber, wg *sync.WaitGroup, csvService *mongoutils.DBService, fm *root.FileMeta) {
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

func writeMetaToMongo(metaService *mongoutils.DBService, fm *root.FileMeta) *root.FileMeta {
	_, inErr := metaService.Collection.InsertOne(metaService.Context, fm)
	if inErr != nil {
		fmt.Println(inErr)
	}
	return fm
}
