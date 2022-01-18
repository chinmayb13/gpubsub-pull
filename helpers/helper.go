package helpers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func WriteResponse(w http.ResponseWriter, contentType string, msg string, count uint64, httpStatus int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(httpStatus)
	response := make(map[string]string)
	response["message"] = msg
	response["acknowledgedMessages"] = strconv.FormatUint(count, 10)
	jsonResp, _ := json.Marshal(response)
	w.Write(jsonResp)
}
