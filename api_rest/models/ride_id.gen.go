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

const rideIDPrefixLength = 1 + 1

var rideIDPrefix = []byte("r_")

type RideID struct {
	raw bunnyid.Raw
}

func (id RideID) Raw() bunnyid.Raw {
	return id.raw
}

func NewRideID() RideID {
	return RideIDFromTime(time.Now())
}

func (id RideID) NextAfter() RideID {
	for i := bunnyid.RawLen - 1; i >= 0; i-- {
		id.raw[i]++
		if id.raw[i] != 0 {
			break
		}
	}
	return id
}

func (id RideID) After(other RideID) bool {
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

func (id RideID) Before(other RideID) bool {
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

func RideIDFromString(id string) (RideID, error) {
	i := &RideID{}
	err := i.UnmarshalText([]byte(id))
	return *i, err
}

func MustRideIDFromString(id string) RideID {
	i := &RideID{}
	err := i.UnmarshalText([]byte(id))
	if err != nil {
		panic(err)
	}
	return *i
}

func RideIDFromTime(t time.Time) RideID {
	var id RideID
	binary.BigEndian.PutUint64(id.raw[:], uint64(t.UnixNano()))
	if _, err := rand.Read(id.raw[6:12]); err != nil {
		panic(errors.Errorf("cannot generate random number: %v;", err))
	}
	return id
}

func (id RideID) String() string {
	text := make([]byte, rideIDPrefixLength+bunnyid.EncodedLen)
	copy(text, rideIDPrefix)
	id.raw.Encode(text[rideIDPrefixLength:])
	return string(text)
}

func (id RideID) MarshalText() ([]byte, error) {
	text := make([]byte, rideIDPrefixLength+bunnyid.EncodedLen)
	copy(text, rideIDPrefix)
	id.raw.Encode(text[rideIDPrefixLength:])
	return text, nil
}

func (id RideID) MarshalJSON() ([]byte, error) {
	text := make([]byte, rideIDPrefixLength+bunnyid.EncodedLen+2)
	text[0] = '"'
	copy(text[1:], rideIDPrefix)
	id.raw.Encode(text[1+rideIDPrefixLength:])
	text[len(text)-1] = '"'
	return text, nil
}

func (id *RideID) UnmarshalText(text []byte) error {
	if len(text) < rideIDPrefixLength {
		return &bunnyid.InvalidIDError{Value: text, Type: "ride_id"}
	}
	if !bytes.Equal(text[:rideIDPrefixLength], rideIDPrefix) {
		parts := strings.Split(string(text), "_")
		if idType, ok := idPrefixes[parts[0]]; ok {
			return &bunnyid.InvalidIDError{Value: text, Type: "ride_id", DetectedType: idType}
		}
		return &bunnyid.InvalidIDError{Value: text, Type: "ride_id"}
	}
	if len(text) != rideIDPrefixLength+bunnyid.EncodedLen {
		return &bunnyid.InvalidIDError{Value: text, Type: "ride_id"}
	}
	text = text[rideIDPrefixLength:]
	if !id.raw.Decode(text) {
		return &bunnyid.InvalidIDError{Value: text, Type: "ride_id"}
	}
	return nil
}

func (id RideID) Time() time.Time {

	var nowBytes [8]byte
	copy(nowBytes[0:6], id.raw[0:6])
	nanos := int64(binary.BigEndian.Uint64(nowBytes[:]))
	return time.Unix(0, nanos).UTC()
}

func (id RideID) Counter() uint64 {
	b := id.raw[6:]

	return uint64(uint64(b[0])<<40 | uint64(b[1])<<32 | uint64(b[2])<<24 | uint64(b[3])<<16 | uint64(b[4])<<8 | uint64(b[5]))
}

func (id RideID) Value() (driver.Value, error) {
	return id.raw[:], nil
}

func (id *RideID) Scan(value interface{}) (err error) {
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
		*id = RideID{}
		return nil
	default:
		return errors.Errorf("xid: scanning unsupported type: %T", value)
	}
}

type NullRideID struct {
	ID    RideID
	Valid bool
}

func NewNullRideID(i RideID, valid bool) NullRideID {
	return NullRideID{
		ID:    i,
		Valid: valid,
	}
}

func NullRideIDFrom(i RideID) NullRideID {
	return NewNullRideID(i, true)
}

func NullRideIDFromPtr(i *RideID) NullRideID {
	if i == nil {
		return NewNullRideID(RideID{}, false)
	}
	return NewNullRideID(*i, true)
}

func (u *NullRideID) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, bunny.NullBytes) {
		u.ID = RideID{}
		u.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &u.ID); err != nil {
		return err
	}

	u.Valid = true
	return nil
}

func (u *NullRideID) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		u.Valid = false
		return nil
	}
	var err error
	res, err := RideIDFromString(string(text))
	u.Valid = err == nil
	if u.Valid {
		u.ID = res
	}
	return err
}

func (u NullRideID) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return bunny.NullBytes, nil
	}
	return u.ID.MarshalJSON()
}

func (u NullRideID) MarshalText() ([]byte, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.ID.MarshalText()
}

func (u *NullRideID) SetValid(n RideID) {
	u.ID = n
	u.Valid = true
}

func (u NullRideID) Ptr() *RideID {
	if !u.Valid {
		return nil
	}
	return &u.ID
}

func (u NullRideID) IsZero() bool {
	return !u.Valid
}

func (u *NullRideID) Scan(value interface{}) error {
	if value == nil {
		u.ID, u.Valid = RideID{}, false
		return nil
	}
	u.Valid = true
	return u.ID.Scan(value)
}

func (u NullRideID) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.ID.Value()
}

func (u NullRideID) String() string {
	if !u.Valid {
		return "<null RideID>"
	}
	return u.ID.String()
}

type RideIDArray []RideID

func (a *RideIDArray) Scan(src interface{}) error {
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

func (a *RideIDArray) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "RideIDArray")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(RideIDArray, len(elems))
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

func (a RideIDArray) Value() (driver.Value, error) {
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
