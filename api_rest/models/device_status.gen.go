package models

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/sqlbunny/runtime/bunny"
	"github.com/sqlbunny/sqlbunny/types/null/convert"
)

type DeviceStatus int32

var DeviceStatuses = struct {
	None            DeviceStatus
	Connected       DeviceStatus
	NotConnected    DeviceStatus
	ConnectionError DeviceStatus
}{
	None:            DeviceStatus(0),
	Connected:       DeviceStatus(1),
	NotConnected:    DeviceStatus(2),
	ConnectionError: DeviceStatus(3),
}

const ()

var deviceStatusValues = map[string]DeviceStatus{
	"None":             DeviceStatus(0),
	"connected":        DeviceStatus(1),
	"not_connected":    DeviceStatus(2),
	"connection_error": DeviceStatus(3),
}

var deviceStatusNames = map[DeviceStatus]string{
	DeviceStatus(0): "None",
	DeviceStatus(1): "connected",
	DeviceStatus(2): "not_connected",
	DeviceStatus(3): "connection_error",
}

func (o DeviceStatus) String() string {
	return deviceStatusNames[o]
}

func DeviceStatusFromString(s string) (DeviceStatus, error) {
	var o DeviceStatus
	err := o.UnmarshalText([]byte(s))
	return o, err
}

func (o DeviceStatus) MarshalText() ([]byte, error) {
	return []byte(o.String()), nil
}

func (o *DeviceStatus) UnmarshalText(text []byte) error {
	val, ok := deviceStatusValues[string(text)]
	if !ok {
		return &bunny.InvalidEnumError{Value: text, Type: "DeviceStatus"}
	}
	*o = val
	return nil
}

type NullDeviceStatus struct {
	DeviceStatus DeviceStatus
	Valid        bool
}

func NewNullDeviceStatus(i DeviceStatus, valid bool) NullDeviceStatus {
	return NullDeviceStatus{
		DeviceStatus: i,
		Valid:        valid,
	}
}

func NullDeviceStatusFrom(i DeviceStatus) NullDeviceStatus {
	return NewNullDeviceStatus(i, true)
}

func NullDeviceStatusFromPtr(i *DeviceStatus) NullDeviceStatus {
	if i == nil {
		var z DeviceStatus
		return NewNullDeviceStatus(z, false)
	}
	return NewNullDeviceStatus(*i, true)
}

func (u *NullDeviceStatus) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, bunny.NullBytes) {
		var z DeviceStatus
		u.DeviceStatus = z
		u.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &u.DeviceStatus); err != nil {
		return err
	}

	u.Valid = true
	return nil
}

func (u *NullDeviceStatus) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		u.Valid = false
		return nil
	}
	var err error
	res, err := DeviceStatusFromString(string(text))
	u.Valid = err == nil
	if u.Valid {
		u.DeviceStatus = res
	}
	return err
}

func (u NullDeviceStatus) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return bunny.NullBytes, nil
	}
	return json.Marshal(u.DeviceStatus)
}

func (u NullDeviceStatus) MarshalText() ([]byte, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.DeviceStatus.MarshalText()
}

func (u *NullDeviceStatus) SetValid(n DeviceStatus) {
	u.DeviceStatus = n
	u.Valid = true
}

func (u NullDeviceStatus) Ptr() *DeviceStatus {
	if !u.Valid {
		return nil
	}
	return &u.DeviceStatus
}

func (u NullDeviceStatus) IsZero() bool {
	return !u.Valid
}

func (u *NullDeviceStatus) Scan(value interface{}) error {
	if value == nil {
		var z DeviceStatus
		u.DeviceStatus, u.Valid = z, false
		return nil
	}
	u.Valid = true
	return convert.ConvertAssign(&u.DeviceStatus, value)
}

func (u NullDeviceStatus) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	return int64(u.DeviceStatus), nil
}

func (u NullDeviceStatus) String() string {
	if !u.Valid {
		return "<null DeviceStatus>"
	}
	return u.DeviceStatus.String()
}
