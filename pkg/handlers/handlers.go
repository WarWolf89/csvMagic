package handlers

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	root ".."
	reader "../reader"
	restutils "../restutils"
	validutils "../validutils"
)

// Uploads a File and send it to be processed
func CsvUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	name := strings.Split(header.Filename, ".")[0]
	fm, err := reader.ProcessCsv(file, name)
	if err != nil {
		restutils.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	restutils.RespondWithJSON(w, http.StatusCreated, fm)
}

func ValidateSingleNumber(w http.ResponseWriter, req *http.Request) {
	procRes := &root.ProcRes{}
	var r io.Reader
	r = req.Body
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)

	b, err := ioutil.ReadAll(tee)
	if err != nil {
		log.Fatal(err)
	}
	num := string(b)
	validutils.CheckAndFixSingleNumber(num, procRes)
	restutils.RespondWithJSON(w, http.StatusCreated, procRes)
}
