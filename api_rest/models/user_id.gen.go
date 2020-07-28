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

const userIDPrefixLength = 3 + 1

var userIDPrefix = []byte("usr_")

type UserID struct {
	raw bunnyid.Raw
}

func (id UserID) Raw() bunnyid.Raw {
	return id.raw
}

func NewUserID() UserID {
	return UserIDFromTime(time.Now())
}

func (id UserID) NextAfter() UserID {
	for i := bunnyid.RawLen - 1; i >= 0; i-- {
		id.raw[i]++
		if id.raw[i] != 0 {
			break
		}
	}
	return id
}

func (id UserID) After(other UserID) bool {
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

func (id UserID) Before(other UserID) bool {
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

func UserIDFromString(id string) (UserID, error) {
	i := &UserID{}
	err := i.UnmarshalText([]byte(id))
	return *i, err
}

func MustUserIDFromString(id string) UserID {
	i := &UserID{}
	err := i.UnmarshalText([]byte(id))
	if err != nil {
		panic(err)
	}
	return *i
}

func UserIDFromTime(t time.Time) UserID {
	var id UserID
	binary.BigEndian.PutUint64(id.raw[:], uint64(t.UnixNano()))
	if _, err := rand.Read(id.raw[6:12]); err != nil {
		panic(errors.Errorf("cannot generate random number: %v;", err))
	}
	return id
}

func (id UserID) String() string {
	text := make([]byte, userIDPrefixLength+bunnyid.EncodedLen)
	copy(text, userIDPrefix)
	id.raw.Encode(text[userIDPrefixLength:])
	return string(text)
}

func (id UserID) MarshalText() ([]byte, error) {
	text := make([]byte, userIDPrefixLength+bunnyid.EncodedLen)
	copy(text, userIDPrefix)
	id.raw.Encode(text[userIDPrefixLength:])
	return text, nil
}

func (id UserID) MarshalJSON() ([]byte, error) {
	text := make([]byte, userIDPrefixLength+bunnyid.EncodedLen+2)
	text[0] = '"'
	copy(text[1:], userIDPrefix)
	id.raw.Encode(text[1+userIDPrefixLength:])
	text[len(text)-1] = '"'
	return text, nil
}

func (id *UserID) UnmarshalText(text []byte) error {
	if len(text) < userIDPrefixLength {
		return &bunnyid.InvalidIDError{Value: text, Type: "user_id"}
	}
	if !bytes.Equal(text[:userIDPrefixLength], userIDPrefix) {
		parts := strings.Split(string(text), "_")
		if idType, ok := idPrefixes[parts[0]]; ok {
			return &bunnyid.InvalidIDError{Value: text, Type: "user_id", DetectedType: idType}
		}
		return &bunnyid.InvalidIDError{Value: text, Type: "user_id"}
	}
	if len(text) != userIDPrefixLength+bunnyid.EncodedLen {
		return &bunnyid.InvalidIDError{Value: text, Type: "user_id"}
	}
	text = text[userIDPrefixLength:]
	if !id.raw.Decode(text) {
		return &bunnyid.InvalidIDError{Value: text, Type: "user_id"}
	}
	return nil
}

func (id UserID) Time() time.Time {

	var nowBytes [8]byte
	copy(nowBytes[0:6], id.raw[0:6])
	nanos := int64(binary.BigEndian.Uint64(nowBytes[:]))
	return time.Unix(0, nanos).UTC()
}

func (id UserID) Counter() uint64 {
	b := id.raw[6:]

	return uint64(uint64(b[0])<<40 | uint64(b[1])<<32 | uint64(b[2])<<24 | uint64(b[3])<<16 | uint64(b[4])<<8 | uint64(b[5]))
}

func (id UserID) Value() (driver.Value, error) {
	return id.raw[:], nil
}

func (id *UserID) Scan(value interface{}) (err error) {
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
		*id = UserID{}
		return nil
	default:
		return errors.Errorf("xid: scanning unsupported type: %T", value)
	}
}

type NullUserID struct {
	ID    UserID
	Valid bool
}

func NewNullUserID(i UserID, valid bool) NullUserID {
	return NullUserID{
		ID:    i,
		Valid: valid,
	}
}

func NullUserIDFrom(i UserID) NullUserID {
	return NewNullUserID(i, true)
}

func NullUserIDFromPtr(i *UserID) NullUserID {
	if i == nil {
		return NewNullUserID(UserID{}, false)
	}
	return NewNullUserID(*i, true)
}

func (u *NullUserID) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, bunny.NullBytes) {
		u.ID = UserID{}
		u.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &u.ID); err != nil {
		return err
	}

	u.Valid = true
	return nil
}

func (u *NullUserID) UnmarshalText(text []byte) error {
	if text == nil || len(text) == 0 {
		u.Valid = false
		return nil
	}
	var err error
	res, err := UserIDFromString(string(text))
	u.Valid = err == nil
	if u.Valid {
		u.ID = res
	}
	return err
}

func (u NullUserID) MarshalJSON() ([]byte, error) {
	if !u.Valid {
		return bunny.NullBytes, nil
	}
	return u.ID.MarshalJSON()
}

func (u NullUserID) MarshalText() ([]byte, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.ID.MarshalText()
}

func (u *NullUserID) SetValid(n UserID) {
	u.ID = n
	u.Valid = true
}

func (u NullUserID) Ptr() *UserID {
	if !u.Valid {
		return nil
	}
	return &u.ID
}

func (u NullUserID) IsZero() bool {
	return !u.Valid
}

func (u *NullUserID) Scan(value interface{}) error {
	if value == nil {
		u.ID, u.Valid = UserID{}, false
		return nil
	}
	u.Valid = true
	return u.ID.Scan(value)
}

func (u NullUserID) Value() (driver.Value, error) {
	if !u.Valid {
		return nil, nil
	}
	return u.ID.Value()
}

func (u NullUserID) String() string {
	if !u.Valid {
		return "<null UserID>"
	}
	return u.ID.String()
}

type UserIDArray []UserID

func (a *UserIDArray) Scan(src interface{}) error {
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

func (a *UserIDArray) scanBytes(src []byte) error {
	elems, err := scanLinearArray(src, []byte{','}, "UserIDArray")
	if err != nil {
		return err
	}
	if *a != nil && len(elems) == 0 {
		*a = (*a)[:0]
	} else {
		b := make(UserIDArray, len(elems))
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

func (a UserIDArray) Value() (driver.Value, error) {
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
