package restutils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func RespondWithFile(w http.ResponseWriter, code int, payload []byte, name string) {
	fileheader := fmt.Sprintf("attachment; filename=%s", name)
	w.Header().Set("Content-Disposition", fileheader)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(payload)
}
