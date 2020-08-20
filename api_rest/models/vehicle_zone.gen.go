package models

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/sqlbunny/sqlbunny/runtime/bunny"
	"github.com/sqlbunny/sqlbunny/types/null/convert"
)

type VehicleZone int32

var VehicleZones = struct {
	None         VehicleZone
	CarrilBici   VehicleZone
	CarrilZona30 VehicleZone
	Acera        VehicleZone
}{
	None:         VehicleZone(0),
	CarrilBici:   VehicleZone(1),
	CarrilZona30: VehicleZone(2),
	Acera:        VehicleZone(3),
}

const ()

var vehicleZoneValues = map[string]VehicleZone{
	"None":          VehicleZone(0),
	"carril_bici":   VehicleZone(1),
	"carril_zona30": VehicleZone(2),
	"acera":         VehicleZone(3),
}

var vehicleZoneNames = map[VehicleZone]string{
	VehicleZone(0): "None",
	VehicleZone(1): "carril_bici",
	VehicleZone(2): "carril_zona30",
	VehicleZone(3): "acera",
}

func (o VehicleZone) String() string {
	return vehicleZoneNames[o]
}

func VehicleZoneFromString(s string) (VehicleZone, error) {
	var o VehicleZone
	err := o.UnmarshalText([]byte(s))
	return o, err
}

func (o VehicleZone) MarshalText() ([]byte, error) {
	return []byte(o.String()), nil
}

func (o *VehicleZone) UnmarshalText(text []byte) error {
	val, ok := vehicleZoneValues[string(text)]
	if !ok {
		return &bunny.InvalidEnumError{Value: text, Type: "VehicleZone"}
	}
	*o = val
	return nil
}

type NullVehicleZone struct {
	VehicleZone VehicleZone
	Valid       bool
}

func NewNullVehicleZone(i VehicleZone, valid bool) NullVehicleZone {
	return NullVehicleZone{
		VehicleZone: i,
		Valid:       valid,
	}
}

func NullVehicleZoneFrom(i VehicleZone) NullVehicleZone {
	return NewNullVehicleZone(i, true)
}

func NullVehicleZoneFromPtr(i *VehicleZone) NullVehicleZone {
	if i == nil {
		var z VehicleZone
		return NewNullVehicleZone(z, false)
	}
	return NewNullVehicleZone(*i, true)
}

func (u *NullVehicleZone) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, bunny.NullBytes) {
		var z VehicleZone
		u.VehicleZone = z
		u.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &u.VehicleZone); err != nil {
		return err
	}

	u.Valid = true
	return nil
}

func (u *NullVehicleZone) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		u.Valid = false
		return nil
	}
	var err error
	res, err := VehicleZoneFromString(string(text))
	u.Valid = err == nil
	if u.Valid {
		u.VehicleZone = res
	}
	return err
}

func (u NullVehicleZone) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return bunny.NullBytes, nil
	}
	return json.Marshal(u.VehicleZone)
}

func (u NullVehicleZone) MarshalText() ([]byte, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.VehicleZone.MarshalText()
}

func (u *NullVehicleZone) SetValid(n VehicleZone) {
	u.VehicleZone = n
	u.Valid = true
}

func (u NullVehicleZone) Ptr() *VehicleZone {
	if !u.Valid {
		return nil
	}
	return &u.VehicleZone
}

func (u NullVehicleZone) IsZero() bool {
	return !u.Valid
}

func (u *NullVehicleZone) Scan(value interface{}) error {
	if value == nil {
		var z VehicleZone
		u.VehicleZone, u.Valid = z, false
		return nil
	}
	u.Valid = true
	return convert.ConvertAssign(&u.VehicleZone, value)
}

func (u NullVehicleZone) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	return int64(u.VehicleZone), nil
}

func (u NullVehicleZone) String() string {
	if !u.Valid {
		return "<null VehicleZone>"
	}
	return u.VehicleZone.String()
}
