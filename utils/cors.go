package utils

import "net/http"

// SetCors settings the CORS headers
func SetCors(w *http.ResponseWriter) {
	// TODO: Just allow all for the discovery stage
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}
