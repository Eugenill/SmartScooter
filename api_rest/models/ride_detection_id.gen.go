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

const rideDetectionIDPrefixLength = 2 + 1

var rideDetectionIDPrefix = []byte("rd_")

type RideDetectionID struct {
	raw bunnyid.Raw
}

func (id RideDetectionID) Raw() bunnyid.Raw {
	return id.raw
}

func NewRideDetectionID() RideDetectionID {
	return RideDetectionIDFromTime(time.Now())
}

func (id RideDetectionID) NextAfter() RideDetectionID {
	for i := bunnyid.RawLen - 1; i >= 0; i-- {
		id.raw[i]++
		if id.raw[i] != 0 {
			break
		}
	}
	return id
}

func (id RideDetectionID) After(other RideDetectionID) bool {
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

func (id RideDetectionID) Before(other RideDetectionID) bool {
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

func RideDetectionIDFromString(id string) (RideDetectionID, error) {
	i := &RideDetectionID{}
	err := i.UnmarshalText([]byte(id))
	return *i, err
}

func MustRideDetectionIDFromString(id string) RideDetectionID {
	i := &RideDetectionID{}
	err := i.UnmarshalText([]byte(id))
	if err != nil {
		panic(err)
	}
	return *i
}

func RideDetectionIDFromTime(t time.Time) RideDetectionID {
	var id RideDetectionID
	binary.BigEndian.PutUint64(id.raw[:], uint64(t.UnixNano()))
	if _, err := rand.Read(id.raw[6:12]); err != nil {
		panic(errors.Errorf("cannot generate random number: %v;", err))
	}
	return id
}

func (id RideDetectionID) String() string {
	text := make([]byte, rideDetectionIDPrefixLength+bunnyid.EncodedLen)
	copy(text, rideDetectionIDPrefix)
	id.raw.Encode(text[rideDetectionIDPrefixLength:])
	return string(text)
}

func (id RideDetectionID) MarshalText() ([]byte, error) {
	text := make([]byte, rideDetectionIDPrefixLength+bunnyid.EncodedLen)
	copy(text, rideDetectionIDPrefix)
	id.raw.Encode(text[rideDetectionIDPrefixLength:])
	return text, nil
}

func (id RideDetectionID) MarshalJSON() ([]byte, error) {
	text := make([]byte, rideDetectionIDPrefixLength+bunnyid.EncodedLen+2)
	text[0] = '"'
	copy(text[1:], rideDetectionIDPrefix)
	id.raw.Encode(text[1+rideDetectionIDPrefixLength:])
	text[len(text)-1] = '"'
	return text, nil
}

func (id *RideDetectionID) UnmarshalText(text []byte) error {
	if len(text) < rideDetectionIDPrefixLength {
		return &bunnyid.InvalidIDError{Value: text, Type: "ride_detection_id"}
	}
	if !bytes.Equal(text[:rideDetectionIDPrefixLength], rideDetectionIDPrefix) {
		parts := strings.Split(string(text), "_")
		if idType, ok := idPrefixes[parts[0]]; ok {
			return &bunnyid.InvalidIDError{Value: text, Type: "ride_detection_id", DetectedType: idType}
		}
		return &bunnyid.InvalidIDError{Value: text, Type: "ride_detection_id"}
	}
	if len(text) != rideDetectionIDPrefixLength+bunnyid.EncodedLen {
		return &bunnyid.InvalidIDError{Value: text, Type: "ride_detection_id"}
	}
	text = text[rideDetectionIDPrefixLength:]
	if !id.raw.Decode(text) {
		return &bunnyid.InvalidIDError{Value: text, Type: "ride_detection_id"}
	}
	return nil
}

func (id RideDetectionID) Time() time.Time {

	var nowBytes [8]byte
	copy(nowBytes[0:6], id.raw[0:6])
	nanos := int64(binary.BigEndian.Uint64(nowBytes[:]))
	return time.Unix(0, nanos).UTC()
}

func (id RideDetectionID) Counter() uint64 {
	b := id.raw[6:]

	return uint64(uint64(b[0])<<40 | uint64(b[1])<<32 | uint64(b[2])<<24 | uint64(b[3])<<16 | uint64(b[4])<<8 | uint64(b[5]))
}

func (id RideDetectionID) Value() (driver.Value, error) {
	return id.raw[:], nil
}

func (id *RideDetectionID) Scan(value interface{}) (err error) {
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
		*id = RideDetectionID{}
		return nil
	default:
		return errors.Errorf("xid: scanning unsupported type: %T", value)
	}
}

type NullRideDetectionID struct {
	ID    RideDetectionID
	Valid bool
}

func NewNullRideDetectionID(i RideDetectionID, valid bool) NullRideDetectionID {
	return NullRideDetectionID{
		ID:    i,
		Valid: valid,
	}
}

func NullRideDetectionIDFrom(i RideDetectionID) NullRideDetectionID {
	return NewNullRideDetectionID(i, true)
}

func NullRideDetectionIDFromPtr(i *RideDetectionID) NullRideDetectionID {
	if i == nil {
		return NewNullRideDetectionID(RideDetectionID{}, false)
	}
	return NewNullRideDetectionID(*i, true)
}

func (u *NullRideDetectionID) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, bunny.NullBytes) {
		u.ID = RideDetectionID{}
		u.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &u.ID); err != nil {
		return err
	}

	u.Valid = true
	return nil
}

func (u *NullRideDetectionID) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		u.Valid = false
		return nil
	}
	var err error
	res, err := RideDetectionIDFromString(string(text))
	u.Valid = err == nil
	if u.Valid {
		u.ID = res
	}
	return err
}

func (u NullRideDetectionID) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return bunny.NullBytes, nil
	}
	return u.ID.MarshalJSON()
}

func (u NullRideDetectionID) MarshalText() ([]byte, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.ID.MarshalText()
}

func (u *NullRideDetectionID) SetValid(n RideDetectionID) {
	u.ID = n
	u.Valid = true
}

func (u NullRideDetectionID) Ptr() *RideDetectionID {
	if !u.Valid {
		return nil
	}
	return &u.ID
}

func (u NullRideDetectionID) IsZero() bool {
	return !u.Valid
}

func (u *NullRideDetectionID) Scan(value interface{}) error {
	if value == nil {
		u.ID, u.Valid = RideDetectionID{}, false
		return nil
	}
	u.Valid = true
	return u.ID.Scan(value)
}

func (u NullRideDetectionID) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.ID.Value()
}

func (u NullRideDetectionID) String() string {
	if !u.Valid {
		return "<null RideDetectionID>"
	}
	return u.ID.String()
}

type RideDetectionIDArray []RideDetectionID

func (a *RideDetectionIDArray) Scan(src interface{}) error {
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

func (a *RideDetectionIDArray) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "RideDetectionIDArray")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(RideDetectionIDArray, len(elems))
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

func (a RideDetectionIDArray) Value() (driver.Value, error) {
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
