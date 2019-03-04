package handlers

import (
	"net/http"
	"strings"

	reader "../reader"
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
	reader.ProcessCsv(file, name)
	return
}
