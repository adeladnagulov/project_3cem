package responses

import (
	"encoding/json"
	"net/http"
)

func SendJSONResp(w http.ResponseWriter, resp map[string]interface{}, code int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(resp)
}
