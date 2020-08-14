package migrations

import "github.com/sqlbunny/sqlbunny/runtime/migration"

func init() {
	Store.Register(&migration.Migration{
		Name:         "00002_05365d",
		Dependencies: []string{"00001_71972d"},
		Operations: []migration.Operation{
			migration.CreateIndexOperation{
				Name:      "user",
				IndexName: "user___login___idx",
				Columns:   []string{"login"},
			},
			migration.CreateIndexOperation{
				Name:      "user",
				IndexName: "user___contact_email___idx",
				Columns:   []string{"contact_email"},
			},
		},
	})
}
