package utils

import (
	"encoding/json"
	"net/http"
)

//RespondWithError method
func RespondWithError(w http.ResponseWriter, code int, payload interface{}) {
	err := ErrorModel{}
	if code == 404 {
		err.Type = "Not found"
		err.Title = "Requested resouce not found"
		RespondWithJSON(w, code, err)
		return
	}
	if code == 400 {
		err.Type = "Bad request"
		err.Title = "Your request parameters didn't validate."
		RespondWithJSON(w, code, err)
		return
	}
	if code == 500 {
		err.Type = "Server errpr"
		err.Title = "Inernal server error"
		RespondWithJSON(w, code, err)
		return
	}
	RespondWithJSON(w, code, map[string]interface{}{"status": code, "message": payload})
}

//RespondWithJSON method
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		response, _ := json.Marshal(payload)
		w.Write(response)
	}
}

// RespondWithValidationError func
func RespondWithValidationError(w http.ResponseWriter, code int, payload []InvalidParams) {
	err := ErrorModel{}
	err.Type = "Bad request"
	err.Title = "Your request parameters didn't validate."
	err.InvalidParams = payload
	RespondWithJSON(w, code, err)
}
