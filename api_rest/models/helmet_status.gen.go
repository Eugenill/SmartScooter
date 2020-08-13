package models

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/sqlbunny/runtime/bunny"
	"github.com/sqlbunny/sqlbunny/types/null/convert"
)

type HelmetStatus int32

var HelmetStatuses = struct {
	Connected       HelmetStatus
	NotConnected    HelmetStatus
	ConnectionError HelmetStatus
}{
	Connected:       HelmetStatus(0),
	NotConnected:    HelmetStatus(1),
	ConnectionError: HelmetStatus(2),
}

const ()

var helmetStatusValues = map[string]HelmetStatus{
	"connected":        HelmetStatus(0),
	"not_connected":    HelmetStatus(1),
	"connection_error": HelmetStatus(2),
}

var helmetStatusNames = map[HelmetStatus]string{
	HelmetStatus(0): "connected",
	HelmetStatus(1): "not_connected",
	HelmetStatus(2): "connection_error",
}

func (o HelmetStatus) String() string {
	return helmetStatusNames[o]
}

func HelmetStatusFromString(s string) (HelmetStatus, error) {
	var o HelmetStatus
	err := o.UnmarshalText([]byte(s))
	return o, err
}

func (o HelmetStatus) MarshalText() ([]byte, error) {
	return []byte(o.String()), nil
}

func (o *HelmetStatus) UnmarshalText(text []byte) error {
	val, ok := helmetStatusValues[string(text)]
	if !ok {
		return &bunny.InvalidEnumError{Value: text, Type: "HelmetStatus"}
	}
	*o = val
	return nil
}

type NullHelmetStatus struct {
	HelmetStatus HelmetStatus
	Valid        bool
}

func NewNullHelmetStatus(i HelmetStatus, valid bool) NullHelmetStatus {
	return NullHelmetStatus{
		HelmetStatus: i,
		Valid:        valid,
	}
}

func NullHelmetStatusFrom(i HelmetStatus) NullHelmetStatus {
	return NewNullHelmetStatus(i, true)
}

func NullHelmetStatusFromPtr(i *HelmetStatus) NullHelmetStatus {
	if i == nil {
		var z HelmetStatus
		return NewNullHelmetStatus(z, false)
	}
	return NewNullHelmetStatus(*i, true)
}

func (u *NullHelmetStatus) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, bunny.NullBytes) {
		var z HelmetStatus
		u.HelmetStatus = z
		u.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &u.HelmetStatus); err != nil {
		return err
	}

	u.Valid = true
	return nil
}

func (u *NullHelmetStatus) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		u.Valid = false
		return nil
	}
	var err error
	res, err := HelmetStatusFromString(string(text))
	u.Valid = err == nil
	if u.Valid {
		u.HelmetStatus = res
	}
	return err
}

func (u NullHelmetStatus) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return bunny.NullBytes, nil
	}
	return json.Marshal(u.HelmetStatus)
}

func (u NullHelmetStatus) MarshalText() ([]byte, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.HelmetStatus.MarshalText()
}

func (u *NullHelmetStatus) SetValid(n HelmetStatus) {
	u.HelmetStatus = n
	u.Valid = true
}

func (u NullHelmetStatus) Ptr() *HelmetStatus {
	if !u.Valid {
		return nil
	}
	return &u.HelmetStatus
}

func (u NullHelmetStatus) IsZero() bool {
	return !u.Valid
}

func (u *NullHelmetStatus) Scan(value interface{}) error {
	if value == nil {
		var z HelmetStatus
		u.HelmetStatus, u.Valid = z, false
		return nil
	}
	u.Valid = true
	return convert.ConvertAssign(&u.HelmetStatus, value)
}

func (u NullHelmetStatus) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	return int64(u.HelmetStatus), nil
}

func (u NullHelmetStatus) String() string {
	if !u.Valid {
		return "<null HelmetStatus>"
	}
	return u.HelmetStatus.String()
}
