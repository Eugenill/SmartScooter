package main

import (
	"github.com/Eugenill/SmartScooter/api_rest/errors"
	"github.com/Eugenill/SmartScooter/api_rest/serve/router"
	"net/http"
)

func main() {
	err := http.ListenAndServe("localhost:1234", router.SetRouter("Ford Mustang"))
	errors.Catch(err)
}
