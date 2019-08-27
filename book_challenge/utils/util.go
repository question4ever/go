package utils

import (
	"encoding/json"
	"net/http"
)

/*
Message(status bool, message string)
JSON message used to give feedback on API funtions
*/
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

/*
Respond(w http.ResponseWriter, data map[string]interface{})
send out the JSON response of the API
*/
func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
