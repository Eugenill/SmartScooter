package reby

import "net/http"

const BearerETSEIB string = "Bearer sess_3quryhvtzpp52va3str1_b851ee78f84c6c4bdebb43e1eb27498397897cb53c8519dc"

func SetHeaders(vehID string, bearer string) http.Header {
	hdr := http.Header{}
	hdr.Set("Content-Type", "application/json")
	hdr.Set("Authorization", bearer)
	hdr.Set("vehicle_id", vehID)
	return hdr
}
