package models

import (
	bytes "bytes"
	rand "crypto/rand"
	driver "database/sql/driver"
	binary "encoding/binary"
	hex "encoding/hex"
	json "encoding/json"
	fmt "fmt"
	bunnyid "github.com/sqlbunny/bunnyid"
	errors "github.com/sqlbunny/errors"
	bunny "github.com/sqlbunny/sqlbunny/runtime/bunny"
	strings "strings"
	time "time"
)

const vehicleIDPrefixLength = 1 + 1

var vehicleIDPrefix = []byte("v_")

type VehicleID struct {
	raw bunnyid.Raw
}

func (id VehicleID) Raw() bunnyid.Raw {
	return id.raw
}

func NewVehicleID() VehicleID {
	return VehicleIDFromTime(time.Now())
}

func (id VehicleID) NextAfter() VehicleID {
	for i := bunnyid.RawLen - 1; i >= 0; i-- {
		id.raw[i]++
		if id.raw[i] != 0 {
			break
		}
	}
	return id
}

func (id VehicleID) After(other VehicleID) bool {
	for i := 0; i < bunnyid.RawLen; i++ {
		if id.raw[i] > other.raw[i] {
			return true
		}
		if id.raw[i] < other.raw[i] {
			return false
		}
	}
	return false
}

func (id VehicleID) Before(other VehicleID) bool {
	for i := 0; i < bunnyid.RawLen; i++ {
		if id.raw[i] < other.raw[i] {
			return true
		}
		if id.raw[i] > other.raw[i] {
			return false
		}
	}
	return false
}

func VehicleIDFromString(id string) (VehicleID, error) {
	i := &VehicleID{}
	err := i.UnmarshalText([]byte(id))
	return *i, err
}

func MustVehicleIDFromString(id string) VehicleID {
	i := &VehicleID{}
	err := i.UnmarshalText([]byte(id))
	if err != nil {
		panic(err)
	}
	return *i
}

func VehicleIDFromTime(t time.Time) VehicleID {
	var id VehicleID
	binary.BigEndian.PutUint64(id.raw[:], uint64(t.UnixNano()))
	if _, err := rand.Read(id.raw[6:12]); err != nil {
		panic(errors.Errorf("cannot generate random number: %v;", err))
	}
	return id
}

func (id VehicleID) String() string {
	text := make([]byte, vehicleIDPrefixLength+bunnyid.EncodedLen)
	copy(text, vehicleIDPrefix)
	id.raw.Encode(text[vehicleIDPrefixLength:])
	return string(text)
}

func (id VehicleID) MarshalText() ([]byte, error) {
	text := make([]byte, vehicleIDPrefixLength+bunnyid.EncodedLen)
	copy(text, vehicleIDPrefix)
	id.raw.Encode(text[vehicleIDPrefixLength:])
	return text, nil
}

func (id VehicleID) MarshalJSON() ([]byte, error) {
	text := make([]byte, vehicleIDPrefixLength+bunnyid.EncodedLen+2)
	text[0] = '"'
	copy(text[1:], vehicleIDPrefix)
	id.raw.Encode(text[1+vehicleIDPrefixLength:])
	text[len(text)-1] = '"'
	return text, nil
}

func (id *VehicleID) UnmarshalText(text []byte) error {
	if len(text) < vehicleIDPrefixLength {
		return &bunnyid.InvalidIDError{Value: text, Type: "vehicle_id"}
	}
	if !bytes.Equal(text[:vehicleIDPrefixLength], vehicleIDPrefix) {
		parts := strings.Split(string(text), "_")
		if idType, ok := idPrefixes[parts[0]]; ok {
			return &bunnyid.InvalidIDError{Value: text, Type: "vehicle_id", DetectedType: idType}
		}
		return &bunnyid.InvalidIDError{Value: text, Type: "vehicle_id"}
	}
	if len(text) != vehicleIDPrefixLength+bunnyid.EncodedLen {
		return &bunnyid.InvalidIDError{Value: text, Type: "vehicle_id"}
	}
	text = text[vehicleIDPrefixLength:]
	if !id.raw.Decode(text) {
		return &bunnyid.InvalidIDError{Value: text, Type: "vehicle_id"}
	}
	return nil
}

func (id VehicleID) Time() time.Time {

	var nowBytes [8]byte
	copy(nowBytes[0:6], id.raw[0:6])
	nanos := int64(binary.BigEndian.Uint64(nowBytes[:]))
	return time.Unix(0, nanos).UTC()
}

func (id VehicleID) Counter() uint64 {
	b := id.raw[6:]

	return uint64(uint64(b[0])<<40 | uint64(b[1])<<32 | uint64(b[2])<<24 | uint64(b[3])<<16 | uint64(b[4])<<8 | uint64(b[5]))
}

func (id VehicleID) Value() (driver.Value, error) {
	return id.raw[:], nil
}

func (id *VehicleID) Scan(value interface{}) (err error) {
	switch val := value.(type) {
	case string:
		return id.UnmarshalText([]byte(val))
	case []byte:
		if len(val) != 12 {
			return errors.Errorf("xid: scanning byte slice invalid length: %d", len(val))
		}
		copy(id.raw[:], val[:])
		return nil
	case nil:
		*id = VehicleID{}
		return nil
	default:
		return errors.Errorf("xid: scanning unsupported type: %T", value)
	}
}

type NullVehicleID struct {
	ID    VehicleID
	Valid bool
}

func NewNullVehicleID(i VehicleID, valid bool) NullVehicleID {
	return NullVehicleID{
		ID:    i,
		Valid: valid,
	}
}

func NullVehicleIDFrom(i VehicleID) NullVehicleID {
	return NewNullVehicleID(i, true)
}

func NullVehicleIDFromPtr(i *VehicleID) NullVehicleID {
	if i == nil {
		return NewNullVehicleID(VehicleID{}, false)
	}
	return NewNullVehicleID(*i, true)
}

func (u *NullVehicleID) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, bunny.NullBytes) {
		u.ID = VehicleID{}
		u.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &u.ID); err != nil {
		return err
	}

	u.Valid = true
	return nil
}

func (u *NullVehicleID) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		u.Valid = false
		return nil
	}
	var err error
	res, err := VehicleIDFromString(string(text))
	u.Valid = err == nil
	if u.Valid {
		u.ID = res
	}
	return err
}

func (u NullVehicleID) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return bunny.NullBytes, nil
	}
	return u.ID.MarshalJSON()
}

func (u NullVehicleID) MarshalText() ([]byte, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.ID.MarshalText()
}

func (u *NullVehicleID) SetValid(n VehicleID) {
	u.ID = n
	u.Valid = true
}

func (u NullVehicleID) Ptr() *VehicleID {
	if !u.Valid {
		return nil
	}
	return &u.ID
}

func (u NullVehicleID) IsZero() bool {
	return !u.Valid
}

func (u *NullVehicleID) Scan(value interface{}) error {
	if value == nil {
		u.ID, u.Valid = VehicleID{}, false
		return nil
	}
	u.Valid = true
	return u.ID.Scan(value)
}

func (u NullVehicleID) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.ID.Value()
}

func (u NullVehicleID) String() string {
	if !u.Valid {
		return "<null VehicleID>"
	}
	return u.ID.String()
}

type VehicleIDArray []VehicleID

func (a *VehicleIDArray) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return a.scanBytes(src)
	case string:
		return a.scanBytes([]byte(src))
	case nil:
		*a = nil
		return nil
	}

	return fmt.Errorf("pq: cannot convert %T to ByteaArray", src)
}

func (a *VehicleIDArray) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "VehicleIDArray")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(VehicleIDArray, len(elems))
		for i, v := range elems {
			bytes, err := parseBytea(v)
			if err != nil {
				return fmt.Errorf("could not parse id array index %d: %s", i, err.Error())
			}
			if len(bytes) != bunnyid.RawLen {
				return fmt.Errorf("could not parse id array index %d: got len %d, expected %d", i, len(bytes), bunnyid.RawLen)
			}
			copy(b[i].raw[:], bytes[:])
		}
		*a = b
	}
	return nil
}

func (a VehicleIDArray) Value() (driver.Value, error) {
	if a == nil {

		return "{}", nil
	}

	if n := len(a); n > 0 {

		size := 1 + 6*n + hex.EncodedLen(bunnyid.RawLen)*len(a)
		b := make([]byte, size)

		for i, s := 0, b; i < n; i++ {
			o := copy(s, `,"\\x`)
			o += hex.Encode(s[o:], a[i].raw[:])
			s[o] = '"'
			s = s[o+1:]
		}

		b[0] = '{'
		b[size-1] = '}'

		return string(b), nil
	}

	return "{}", nil
}
