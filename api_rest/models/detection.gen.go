package models

import (
	_import00 "github.com/sqlbunny/sqlbunny/types/null"
	_import01 "time"
)

import (
	"bytes"
	"encoding/json"

	"github.com/sqlbunny/sqlbunny/runtime/bunny"
)

type Detection struct {
	TrafficLight _import00.String `bunny:"traffic_light" json:"traffic_light" `
	Obstacle     _import00.String `json:"obstacle" bunny:"obstacle" `
	Location     Point            `bunny:"location__,bind" json:"location" `
	DetectedAt   _import01.Time   `bunny:"detected_at" json:"detected_at" `
}
type NullDetection struct {
	Detection Detection
	Valid     bool
}

func NewNullDetection(s Detection, valid bool) NullDetection {
	return NullDetection{
		Detection: s,
		Valid:     valid,
	}
}

func NullDetectionFrom(s Detection) NullDetection {
	return NewNullDetection(s, true)
}

func NullDetectionFromPtr(s *Detection) NullDetection {
	if s == nil {
		return NewNullDetection(Detection{}, false)
	}
	return NewNullDetection(*s, true)
}

func (u *NullDetection) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, bunny.NullBytes) {
		u.Detection = Detection{}
		u.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &u.Detection); err != nil {
		return err
	}

	u.Valid = true
	return nil
}

func (u NullDetection) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return bunny.NullBytes, nil
	}
	return json.Marshal(u.Detection)
}

func (u *NullDetection) SetValid(n Detection) {
	u.Detection = n
	u.Valid = true
}

func (u NullDetection) Ptr() *Detection {
	if !u.Valid {
		return nil
	}
	return &u.Detection
}

func (u NullDetection) IsZero() bool {
	return !u.Valid
}
