package rest

import (
	"encoding/json"
	"io"
	"net/http"
)

func unmarshalJSON(val interface{}, r io.Reader) error {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	if err := dec.Decode(val); err != nil {
		return err
	}
	return nil
}

func UnmarshalJSONRequest(val interface{}, r *http.Request) error {
	if sk, ok := r.Body.(io.Seeker); ok {
		_, err := sk.Seek(0, io.SeekStart)
		if err != nil {
			panic(err) // Should never happen
		}
	}
	return unmarshalJSON(val, r.Body)
}
