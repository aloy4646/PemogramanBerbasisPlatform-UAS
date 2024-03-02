package controller

import (
	"encoding/json"
	"kuis1/model"
	"net/http"
)

func responseFromRowsAffected(w http.ResponseWriter, rowsAffected int64) {
	if rowsAffected > 0 {
		successResponseMessage(w)
	} else {
		errorResponseMessage(w, 407, "Failed, 0 rows affected")
	}
}

func successResponseMessage(w http.ResponseWriter) {
	var response model.SuccessResponse
	response.Status = 200
	response.Message = "Success"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func errorResponseMessage(w http.ResponseWriter, status int, message string) {
	var response model.ErrorResponse
	response.Status = status
	response.Message = message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendSuccessResponseWithData(w http.ResponseWriter, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	var response model.SuccessResponseWithData
	response.Status = http.StatusOK
	response.Data = value
	json.NewEncoder(w).Encode(response)
}

func sendUnAuthorizedResponse(w http.ResponseWriter) {
	var response model.ErrorResponse
	response.Status = 401
	response.Message = "UnAuthorized Access"

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
