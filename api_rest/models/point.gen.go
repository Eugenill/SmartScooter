package models

import (
	"bytes"
	"encoding/json"

	"github.com/sqlbunny/sqlbunny/runtime/bunny"
)

type Point struct {
	Latitude  float64 `bunny:"latitude" json:"latitude" `
	Longitude float64 `json:"longitude" bunny:"longitude" `
	Accuracy  float64 `bunny:"accuracy" json:"accuracy" `
}
type NullPoint struct {
	Point Point
	Valid bool
}

func NewNullPoint(s Point, valid bool) NullPoint {
	return NullPoint{
		Point: s,
		Valid: valid,
	}
}

func NullPointFrom(s Point) NullPoint {
	return NewNullPoint(s, true)
}

func NullPointFromPtr(s *Point) NullPoint {
	if s == nil {
		return NewNullPoint(Point{}, false)
	}
	return NewNullPoint(*s, true)
}

func (u *NullPoint) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, bunny.NullBytes) {
		u.Point = Point{}
		u.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &u.Point); err != nil {
		return err
	}

	u.Valid = true
	return nil
}

func (u NullPoint) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return bunny.NullBytes, nil
	}
	return json.Marshal(u.Point)
}

func (u *NullPoint) SetValid(n Point) {
	u.Point = n
	u.Valid = true
}

func (u NullPoint) Ptr() *Point {
	if !u.Valid {
		return nil
	}
	return &u.Point
}

func (u NullPoint) IsZero() bool {
	return !u.Valid
}
