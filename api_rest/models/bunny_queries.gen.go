package models

import (
	"github.com/sqlbunny/sqlbunny/runtime/qm"
	"github.com/sqlbunny/sqlbunny/runtime/queries"
)

var dialect = queries.Dialect{
	LQ:                0x22,
	RQ:                0x22,
	IndexPlaceholders: true,
	UseTopClause:      false,
}

func NewQuery(mods ...qm.QueryMod) *queries.Query {
	q := &queries.Query{}
	queries.SetDialect(q, &dialect)
	qm.Apply(q, mods...)

	return q
}
