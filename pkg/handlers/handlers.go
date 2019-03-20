package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	reader "../reader"
	restutils "../restutils"
	validutils "../validutils"
	"github.com/gorilla/mux"
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
	resData, ferr := json.MarshalIndent(fm, "", " ")
	if ferr != nil {
		restutils.RespondWithError(w, http.StatusInternalServerError, ferr.Error())
	}

	restutils.RespondWithFile(w, http.StatusCreated, resData, name)
}

func ValidateSingleNumber(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	fmt.Println(vars)
	num := vars["num"]
	pr := validutils.CheckAndFixSingleNumber(num)
	restutils.RespondWithJSON(w, http.StatusOK, pr)
}

func GetFileInfo(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	fid := vars["fid"]
	fm, err := reader.GetFileByID(fid)
	if err != nil {
		restutils.RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	restutils.RespondWithJSON(w, http.StatusOK, fm)
}
