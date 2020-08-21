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

const iotDeviceIDPrefixLength = 3 + 1

var iotDeviceIDPrefix = []byte("iot_")

type IotDeviceID struct {
	raw bunnyid.Raw
}

func (id IotDeviceID) Raw() bunnyid.Raw {
	return id.raw
}

func NewIotDeviceID() IotDeviceID {
	return IotDeviceIDFromTime(time.Now())
}

func (id IotDeviceID) NextAfter() IotDeviceID {
	for i := bunnyid.RawLen - 1; i >= 0; i-- {
		id.raw[i]++
		if id.raw[i] != 0 {
			break
		}
	}
	return id
}

func (id IotDeviceID) After(other IotDeviceID) bool {
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

func (id IotDeviceID) Before(other IotDeviceID) bool {
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

func IotDeviceIDFromString(id string) (IotDeviceID, error) {
	i := &IotDeviceID{}
	err := i.UnmarshalText([]byte(id))
	return *i, err
}

func MustIotDeviceIDFromString(id string) IotDeviceID {
	i := &IotDeviceID{}
	err := i.UnmarshalText([]byte(id))
	if err != nil {
		panic(err)
	}
	return *i
}

func IotDeviceIDFromTime(t time.Time) IotDeviceID {
	var id IotDeviceID
	binary.BigEndian.PutUint64(id.raw[:], uint64(t.UnixNano()))
	if _, err := rand.Read(id.raw[6:12]); err != nil {
		panic(errors.Errorf("cannot generate random number: %v;", err))
	}
	return id
}

func (id IotDeviceID) String() string {
	text := make([]byte, iotDeviceIDPrefixLength+bunnyid.EncodedLen)
	copy(text, iotDeviceIDPrefix)
	id.raw.Encode(text[iotDeviceIDPrefixLength:])
	return string(text)
}

func (id IotDeviceID) MarshalText() ([]byte, error) {
	text := make([]byte, iotDeviceIDPrefixLength+bunnyid.EncodedLen)
	copy(text, iotDeviceIDPrefix)
	id.raw.Encode(text[iotDeviceIDPrefixLength:])
	return text, nil
}

func (id IotDeviceID) MarshalJSON() ([]byte, error) {
	text := make([]byte, iotDeviceIDPrefixLength+bunnyid.EncodedLen+2)
	text[0] = '"'
	copy(text[1:], iotDeviceIDPrefix)
	id.raw.Encode(text[1+iotDeviceIDPrefixLength:])
	text[len(text)-1] = '"'
	return text, nil
}

func (id *IotDeviceID) UnmarshalText(text []byte) error {
	if len(text) < iotDeviceIDPrefixLength {
		return &bunnyid.InvalidIDError{Value: text, Type: "iot_device_id"}
	}
	if !bytes.Equal(text[:iotDeviceIDPrefixLength], iotDeviceIDPrefix) {
		parts := strings.Split(string(text), "_")
		if idType, ok := idPrefixes[parts[0]]; ok {
			return &bunnyid.InvalidIDError{Value: text, Type: "iot_device_id", DetectedType: idType}
		}
		return &bunnyid.InvalidIDError{Value: text, Type: "iot_device_id"}
	}
	if len(text) != iotDeviceIDPrefixLength+bunnyid.EncodedLen {
		return &bunnyid.InvalidIDError{Value: text, Type: "iot_device_id"}
	}
	text = text[iotDeviceIDPrefixLength:]
	if !id.raw.Decode(text) {
		return &bunnyid.InvalidIDError{Value: text, Type: "iot_device_id"}
	}
	return nil
}

func (id IotDeviceID) Time() time.Time {

	var nowBytes [8]byte
	copy(nowBytes[0:6], id.raw[0:6])
	nanos := int64(binary.BigEndian.Uint64(nowBytes[:]))
	return time.Unix(0, nanos).UTC()
}

func (id IotDeviceID) Counter() uint64 {
	b := id.raw[6:]

	return uint64(uint64(b[0])<<40 | uint64(b[1])<<32 | uint64(b[2])<<24 | uint64(b[3])<<16 | uint64(b[4])<<8 | uint64(b[5]))
}

func (id IotDeviceID) Value() (driver.Value, error) {
	return id.raw[:], nil
}

func (id *IotDeviceID) Scan(value interface{}) (err error) {
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
		*id = IotDeviceID{}
		return nil
	default:
		return errors.Errorf("xid: scanning unsupported type: %T", value)
	}
}

type NullIotDeviceID struct {
	ID    IotDeviceID
	Valid bool
}

func NewNullIotDeviceID(i IotDeviceID, valid bool) NullIotDeviceID {
	return NullIotDeviceID{
		ID:    i,
		Valid: valid,
	}
}

func NullIotDeviceIDFrom(i IotDeviceID) NullIotDeviceID {
	return NewNullIotDeviceID(i, true)
}

func NullIotDeviceIDFromPtr(i *IotDeviceID) NullIotDeviceID {
	if i == nil {
		return NewNullIotDeviceID(IotDeviceID{}, false)
	}
	return NewNullIotDeviceID(*i, true)
}

func (u *NullIotDeviceID) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, bunny.NullBytes) {
		u.ID = IotDeviceID{}
		u.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &u.ID); err != nil {
		return err
	}

	u.Valid = true
	return nil
}

func (u *NullIotDeviceID) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		u.Valid = false
		return nil
	}
	var err error
	res, err := IotDeviceIDFromString(string(text))
	u.Valid = err == nil
	if u.Valid {
		u.ID = res
	}
	return err
}

func (u NullIotDeviceID) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return bunny.NullBytes, nil
	}
	return u.ID.MarshalJSON()
}

func (u NullIotDeviceID) MarshalText() ([]byte, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.ID.MarshalText()
}

func (u *NullIotDeviceID) SetValid(n IotDeviceID) {
	u.ID = n
	u.Valid = true
}

func (u NullIotDeviceID) Ptr() *IotDeviceID {
	if !u.Valid {
		return nil
	}
	return &u.ID
}

func (u NullIotDeviceID) IsZero() bool {
	return !u.Valid
}

func (u *NullIotDeviceID) Scan(value interface{}) error {
	if value == nil {
		u.ID, u.Valid = IotDeviceID{}, false
		return nil
	}
	u.Valid = true
	return u.ID.Scan(value)
}

func (u NullIotDeviceID) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.ID.Value()
}

func (u NullIotDeviceID) String() string {
	if !u.Valid {
		return "<null IotDeviceID>"
	}
	return u.ID.String()
}

type IotDeviceIDArray []IotDeviceID

func (a *IotDeviceIDArray) Scan(src interface{}) error {
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

func (a *IotDeviceIDArray) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "IotDeviceIDArray")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(IotDeviceIDArray, len(elems))
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

func (a IotDeviceIDArray) Value() (driver.Value, error) {
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
