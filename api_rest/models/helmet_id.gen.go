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

const helmetIDPrefixLength = 1 + 1

var helmetIDPrefix = []byte("h_")

type HelmetID struct {
	raw bunnyid.Raw
}

func (id HelmetID) Raw() bunnyid.Raw {
	return id.raw
}

func NewHelmetID() HelmetID {
	return HelmetIDFromTime(time.Now())
}

func (id HelmetID) NextAfter() HelmetID {
	for i := bunnyid.RawLen - 1; i >= 0; i-- {
		id.raw[i]++
		if id.raw[i] != 0 {
			break
		}
	}
	return id
}

func (id HelmetID) After(other HelmetID) bool {
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

func (id HelmetID) Before(other HelmetID) bool {
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

func HelmetIDFromString(id string) (HelmetID, error) {
	i := &HelmetID{}
	err := i.UnmarshalText([]byte(id))
	return *i, err
}

func MustHelmetIDFromString(id string) HelmetID {
	i := &HelmetID{}
	err := i.UnmarshalText([]byte(id))
	if err != nil {
		panic(err)
	}
	return *i
}

func HelmetIDFromTime(t time.Time) HelmetID {
	var id HelmetID
	binary.BigEndian.PutUint64(id.raw[:], uint64(t.UnixNano()))
	if _, err := rand.Read(id.raw[6:12]); err != nil {
		panic(errors.Errorf("cannot generate random number: %v;", err))
	}
	return id
}

func (id HelmetID) String() string {
	text := make([]byte, helmetIDPrefixLength+bunnyid.EncodedLen)
	copy(text, helmetIDPrefix)
	id.raw.Encode(text[helmetIDPrefixLength:])
	return string(text)
}

func (id HelmetID) MarshalText() ([]byte, error) {
	text := make([]byte, helmetIDPrefixLength+bunnyid.EncodedLen)
	copy(text, helmetIDPrefix)
	id.raw.Encode(text[helmetIDPrefixLength:])
	return text, nil
}

func (id HelmetID) MarshalJSON() ([]byte, error) {
	text := make([]byte, helmetIDPrefixLength+bunnyid.EncodedLen+2)
	text[0] = '"'
	copy(text[1:], helmetIDPrefix)
	id.raw.Encode(text[1+helmetIDPrefixLength:])
	text[len(text)-1] = '"'
	return text, nil
}

func (id *HelmetID) UnmarshalText(text []byte) error {
	if len(text) < helmetIDPrefixLength {
		return &bunnyid.InvalidIDError{Value: text, Type: "helmet_id"}
	}
	if !bytes.Equal(text[:helmetIDPrefixLength], helmetIDPrefix) {
		parts := strings.Split(string(text), "_")
		if idType, ok := idPrefixes[parts[0]]; ok {
			return &bunnyid.InvalidIDError{Value: text, Type: "helmet_id", DetectedType: idType}
		}
		return &bunnyid.InvalidIDError{Value: text, Type: "helmet_id"}
	}
	if len(text) != helmetIDPrefixLength+bunnyid.EncodedLen {
		return &bunnyid.InvalidIDError{Value: text, Type: "helmet_id"}
	}
	text = text[helmetIDPrefixLength:]
	if !id.raw.Decode(text) {
		return &bunnyid.InvalidIDError{Value: text, Type: "helmet_id"}
	}
	return nil
}

func (id HelmetID) Time() time.Time {

	var nowBytes [8]byte
	copy(nowBytes[0:6], id.raw[0:6])
	nanos := int64(binary.BigEndian.Uint64(nowBytes[:]))
	return time.Unix(0, nanos).UTC()
}

func (id HelmetID) Counter() uint64 {
	b := id.raw[6:]

	return uint64(uint64(b[0])<<40 | uint64(b[1])<<32 | uint64(b[2])<<24 | uint64(b[3])<<16 | uint64(b[4])<<8 | uint64(b[5]))
}

func (id HelmetID) Value() (driver.Value, error) {
	return id.raw[:], nil
}

func (id *HelmetID) Scan(value interface{}) (err error) {
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
		*id = HelmetID{}
		return nil
	default:
		return errors.Errorf("xid: scanning unsupported type: %T", value)
	}
}

type NullHelmetID struct {
	ID    HelmetID
	Valid bool
}

func NewNullHelmetID(i HelmetID, valid bool) NullHelmetID {
	return NullHelmetID{
		ID:    i,
		Valid: valid,
	}
}

func NullHelmetIDFrom(i HelmetID) NullHelmetID {
	return NewNullHelmetID(i, true)
}

func NullHelmetIDFromPtr(i *HelmetID) NullHelmetID {
	if i == nil {
		return NewNullHelmetID(HelmetID{}, false)
	}
	return NewNullHelmetID(*i, true)
}

func (u *NullHelmetID) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, bunny.NullBytes) {
		u.ID = HelmetID{}
		u.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &u.ID); err != nil {
		return err
	}

	u.Valid = true
	return nil
}

func (u *NullHelmetID) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		u.Valid = false
		return nil
	}
	var err error
	res, err := HelmetIDFromString(string(text))
	u.Valid = err == nil
	if u.Valid {
		u.ID = res
	}
	return err
}

func (u NullHelmetID) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return bunny.NullBytes, nil
	}
	return u.ID.MarshalJSON()
}

func (u NullHelmetID) MarshalText() ([]byte, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.ID.MarshalText()
}

func (u *NullHelmetID) SetValid(n HelmetID) {
	u.ID = n
	u.Valid = true
}

func (u NullHelmetID) Ptr() *HelmetID {
	if !u.Valid {
		return nil
	}
	return &u.ID
}

func (u NullHelmetID) IsZero() bool {
	return !u.Valid
}

func (u *NullHelmetID) Scan(value interface{}) error {
	if value == nil {
		u.ID, u.Valid = HelmetID{}, false
		return nil
	}
	u.Valid = true
	return u.ID.Scan(value)
}

func (u NullHelmetID) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.ID.Value()
}

func (u NullHelmetID) String() string {
	if !u.Valid {
		return "<null HelmetID>"
	}
	return u.ID.String()
}

type HelmetIDArray []HelmetID

func (a *HelmetIDArray) Scan(src interface{}) error {
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

func (a *HelmetIDArray) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "HelmetIDArray")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(HelmetIDArray, len(elems))
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

func (a HelmetIDArray) Value() (driver.Value, error) {
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
