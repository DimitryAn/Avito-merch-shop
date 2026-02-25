package responser

import (
	"encoding/json"
	"net/http"
)

func SendHttpError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"errors": message})
}
