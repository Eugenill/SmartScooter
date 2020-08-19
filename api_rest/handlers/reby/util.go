package reby

import "net/http"

func SetHeaders(vehID string, bearer string) http.Header {
	hdr := http.Header{}
	hdr.Set("Content-Type", "application/json")
	hdr.Set("Authorization", bearer)
	hdr.Set("vehicle_id", vehID)
	return hdr
}
