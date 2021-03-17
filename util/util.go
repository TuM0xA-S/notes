package util

import (
	"encoding/json"
	"net/http"
)

// Message constructs message object
func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

// Respond responds with json encoding
func Respond(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}
