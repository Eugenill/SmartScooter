package models

import (
	"github.com/sqlbunny/sqlbunny/runtime/queries"
	"github.com/sqlbunny/sqlbunny/runtime/strmangle"
)

type M map[string]interface{}

type insertCache struct {
	query        string
	valueMapping []queries.MappedField
}

type updateCache struct {
	query        string
	valueMapping []queries.MappedField
}

func makeCacheKey(wl []string) string {
	buf := strmangle.GetBuffer()

	for _, w := range wl {
		buf.WriteString(w)
		buf.WriteByte(',')
	}

	str := buf.String()
	strmangle.PutBuffer(buf)
	return str
}
