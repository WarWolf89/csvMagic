package reader

import (
	"../mongoutils"
	"os"
	"testing"
)

const(
	collTest = "test"
)
var (
	testService = mongoutils.NewDBService(addr, db, collTest)
)

func TestReader(t *testing.T) {
	f, err := os.Open("../../resources/csv/test.csv")
	if err != nil {
		panic(err)
	}
	fmTest, err := ProcessCsv(f,"test")
	if (fmTest.Counters["valid"] != 1 || fmTest.Counters["fixed"]  != 1 || fmTest.Counters["unfixable"]  != 1 || fmTest.Counters["processed"] != 3) {
		t.Errorf("Incorrect values have been produced by reader, instead got values %v", fmTest.Counters)
	}
	}