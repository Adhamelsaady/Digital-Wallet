package api

import (
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func HandleHealth (responseWrite http.ResponseWriter , request *http.Request){
	responseWrite.Header().Set("content-type" , "application/json")
	responseWrite.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWrite).Encode(HealthResponse{Status: "ok"})
}
