package migrations

import "github.com/sqlbunny/sqlbunny/runtime/migration"

func init() {
	Store.Register(&migration.Migration{
		Name:         "00005_0479af",
		Dependencies: []string{"00004_f943a0"},
		Operations: []migration.Operation{
			migration.AlterTableOperation{
				Name: "helmet",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableAddColumn{
						Name:     "name",
						Type:     "text",
						Default:  "''",
						Nullable: false,
					},
				},
			},
		},
	})
}
