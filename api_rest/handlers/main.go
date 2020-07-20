package handlers

import (
	"github.com/Eugenill/SmartScooter/api_rest/errors"
	"net/http"
)

func GetVehicles(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(name))
		errors.Catch(err)
	}
}
