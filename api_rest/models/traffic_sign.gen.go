package models

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/sqlbunny/runtime/bunny"
	"github.com/sqlbunny/sqlbunny/types/null/convert"
)

type TrafficSign int32

var TrafficSigns = struct {
	Stop   TrafficSign
	Yield  TrafficSign
	Vel10  TrafficSign
	Vel20  TrafficSign
	Vel30  TrafficSign
	People TrafficSign
}{
	Stop:   TrafficSign(0),
	Yield:  TrafficSign(1),
	Vel10:  TrafficSign(2),
	Vel20:  TrafficSign(3),
	Vel30:  TrafficSign(4),
	People: TrafficSign(5),
}

const ()

var trafficSignValues = map[string]TrafficSign{
	"stop":   TrafficSign(0),
	"yield":  TrafficSign(1),
	"vel_10": TrafficSign(2),
	"vel_20": TrafficSign(3),
	"vel_30": TrafficSign(4),
	"people": TrafficSign(5),
}

var trafficSignNames = map[TrafficSign]string{
	TrafficSign(0): "stop",
	TrafficSign(1): "yield",
	TrafficSign(2): "vel_10",
	TrafficSign(3): "vel_20",
	TrafficSign(4): "vel_30",
	TrafficSign(5): "people",
}

func (o TrafficSign) String() string {
	return trafficSignNames[o]
}

func TrafficSignFromString(s string) (TrafficSign, error) {
	var o TrafficSign
	err := o.UnmarshalText([]byte(s))
	return o, err
}

func (o TrafficSign) MarshalText() ([]byte, error) {
	return []byte(o.String()), nil
}

func (o *TrafficSign) UnmarshalText(text []byte) error {
	val, ok := trafficSignValues[string(text)]
	if !ok {
		return &bunny.InvalidEnumError{Value: text, Type: "TrafficSign"}
	}
	*o = val
	return nil
}

type NullTrafficSign struct {
	TrafficSign TrafficSign
	Valid       bool
}

func NewNullTrafficSign(i TrafficSign, valid bool) NullTrafficSign {
	return NullTrafficSign{
		TrafficSign: i,
		Valid:       valid,
	}
}

func NullTrafficSignFrom(i TrafficSign) NullTrafficSign {
	return NewNullTrafficSign(i, true)
}

func NullTrafficSignFromPtr(i *TrafficSign) NullTrafficSign {
	if i == nil {
		var z TrafficSign
		return NewNullTrafficSign(z, false)
	}
	return NewNullTrafficSign(*i, true)
}

func (u *NullTrafficSign) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, bunny.NullBytes) {
		var z TrafficSign
		u.TrafficSign = z
		u.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &u.TrafficSign); err != nil {
		return err
	}

	u.Valid = true
	return nil
}

func (u *NullTrafficSign) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		u.Valid = false
		return nil
	}
	var err error
	res, err := TrafficSignFromString(string(text))
	u.Valid = err == nil
	if u.Valid {
		u.TrafficSign = res
	}
	return err
}

func (u NullTrafficSign) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return bunny.NullBytes, nil
	}
	return json.Marshal(u.TrafficSign)
}

func (u NullTrafficSign) MarshalText() ([]byte, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.TrafficSign.MarshalText()
}

func (u *NullTrafficSign) SetValid(n TrafficSign) {
	u.TrafficSign = n
	u.Valid = true
}

func (u NullTrafficSign) Ptr() *TrafficSign {
	if !u.Valid {
		return nil
	}
	return &u.TrafficSign
}

func (u NullTrafficSign) IsZero() bool {
	return !u.Valid
}

func (u *NullTrafficSign) Scan(value interface{}) error {
	if value == nil {
		var z TrafficSign
		u.TrafficSign, u.Valid = z, false
		return nil
	}
	u.Valid = true
	return convert.ConvertAssign(&u.TrafficSign, value)
}

func (u NullTrafficSign) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	return int64(u.TrafficSign), nil
}

func (u NullTrafficSign) String() string {
	if !u.Valid {
		return "<null TrafficSign>"
	}
	return u.TrafficSign.String()
}
