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

const pathIDPrefixLength = 2 + 1

var pathIDPrefix = []byte("pa_")

type PathID struct {
	raw bunnyid.Raw
}

func (id PathID) Raw() bunnyid.Raw {
	return id.raw
}

func NewPathID() PathID {
	return PathIDFromTime(time.Now())
}

func (id PathID) NextAfter() PathID {
	for i := bunnyid.RawLen - 1; i >= 0; i-- {
		id.raw[i]++
		if id.raw[i] != 0 {
			break
		}
	}
	return id
}

func (id PathID) After(other PathID) bool {
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

func (id PathID) Before(other PathID) bool {
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

func PathIDFromString(id string) (PathID, error) {
	i := &PathID{}
	err := i.UnmarshalText([]byte(id))
	return *i, err
}

func MustPathIDFromString(id string) PathID {
	i := &PathID{}
	err := i.UnmarshalText([]byte(id))
	if err != nil {
		panic(err)
	}
	return *i
}

func PathIDFromTime(t time.Time) PathID {
	var id PathID
	binary.BigEndian.PutUint64(id.raw[:], uint64(t.UnixNano()))
	if _, err := rand.Read(id.raw[6:12]); err != nil {
		panic(errors.Errorf("cannot generate random number: %v;", err))
	}
	return id
}

func (id PathID) String() string {
	text := make([]byte, pathIDPrefixLength+bunnyid.EncodedLen)
	copy(text, pathIDPrefix)
	id.raw.Encode(text[pathIDPrefixLength:])
	return string(text)
}

func (id PathID) MarshalText() ([]byte, error) {
	text := make([]byte, pathIDPrefixLength+bunnyid.EncodedLen)
	copy(text, pathIDPrefix)
	id.raw.Encode(text[pathIDPrefixLength:])
	return text, nil
}

func (id PathID) MarshalJSON() ([]byte, error) {
	text := make([]byte, pathIDPrefixLength+bunnyid.EncodedLen+2)
	text[0] = '"'
	copy(text[1:], pathIDPrefix)
	id.raw.Encode(text[1+pathIDPrefixLength:])
	text[len(text)-1] = '"'
	return text, nil
}

func (id *PathID) UnmarshalText(text []byte) error {
	if len(text) < pathIDPrefixLength {
		return &bunnyid.InvalidIDError{Value: text, Type: "path_id"}
	}
	if !bytes.Equal(text[:pathIDPrefixLength], pathIDPrefix) {
		parts := strings.Split(string(text), "_")
		if idType, ok := idPrefixes[parts[0]]; ok {
			return &bunnyid.InvalidIDError{Value: text, Type: "path_id", DetectedType: idType}
		}
		return &bunnyid.InvalidIDError{Value: text, Type: "path_id"}
	}
	if len(text) != pathIDPrefixLength+bunnyid.EncodedLen {
		return &bunnyid.InvalidIDError{Value: text, Type: "path_id"}
	}
	text = text[pathIDPrefixLength:]
	if !id.raw.Decode(text) {
		return &bunnyid.InvalidIDError{Value: text, Type: "path_id"}
	}
	return nil
}

func (id PathID) Time() time.Time {

	var nowBytes [8]byte
	copy(nowBytes[0:6], id.raw[0:6])
	nanos := int64(binary.BigEndian.Uint64(nowBytes[:]))
	return time.Unix(0, nanos).UTC()
}

func (id PathID) Counter() uint64 {
	b := id.raw[6:]

	return uint64(uint64(b[0])<<40 | uint64(b[1])<<32 | uint64(b[2])<<24 | uint64(b[3])<<16 | uint64(b[4])<<8 | uint64(b[5]))
}

func (id PathID) Value() (driver.Value, error) {
	return id.raw[:], nil
}

func (id *PathID) Scan(value interface{}) (err error) {
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
		*id = PathID{}
		return nil
	default:
		return errors.Errorf("xid: scanning unsupported type: %T", value)
	}
}

type NullPathID struct {
	ID    PathID
	Valid bool
}

func NewNullPathID(i PathID, valid bool) NullPathID {
	return NullPathID{
		ID:    i,
		Valid: valid,
	}
}

func NullPathIDFrom(i PathID) NullPathID {
	return NewNullPathID(i, true)
}

func NullPathIDFromPtr(i *PathID) NullPathID {
	if i == nil {
		return NewNullPathID(PathID{}, false)
	}
	return NewNullPathID(*i, true)
}

func (u *NullPathID) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, bunny.NullBytes) {
		u.ID = PathID{}
		u.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &u.ID); err != nil {
		return err
	}

	u.Valid = true
	return nil
}

func (u *NullPathID) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		u.Valid = false
		return nil
	}
	var err error
	res, err := PathIDFromString(string(text))
	u.Valid = err == nil
	if u.Valid {
		u.ID = res
	}
	return err
}

func (u NullPathID) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return bunny.NullBytes, nil
	}
	return u.ID.MarshalJSON()
}

func (u NullPathID) MarshalText() ([]byte, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.ID.MarshalText()
}

func (u *NullPathID) SetValid(n PathID) {
	u.ID = n
	u.Valid = true
}

func (u NullPathID) Ptr() *PathID {
	if !u.Valid {
		return nil
	}
	return &u.ID
}

func (u NullPathID) IsZero() bool {
	return !u.Valid
}

func (u *NullPathID) Scan(value interface{}) error {
	if value == nil {
		u.ID, u.Valid = PathID{}, false
		return nil
	}
	u.Valid = true
	return u.ID.Scan(value)
}

func (u NullPathID) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.ID.Value()
}

func (u NullPathID) String() string {
	if !u.Valid {
		return "<null PathID>"
	}
	return u.ID.String()
}

type PathIDArray []PathID

func (a *PathIDArray) Scan(src interface{}) error {
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

func (a *PathIDArray) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "PathIDArray")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(PathIDArray, len(elems))
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

func (a PathIDArray) Value() (driver.Value, error) {
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
