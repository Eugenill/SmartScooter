package handlers

import (
	"github.com/Eugenill/SmartScooter/api_rest/pkg/errors"
	"log"
	"net/http"
)

func GetVehicles(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r)
		_, err := w.Write([]byte(name))
		errors.Catch(err)
	}
}
