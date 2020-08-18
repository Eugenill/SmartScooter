package migrations

import "github.com/sqlbunny/sqlbunny/runtime/migration"

func init() {
	Store.Register(&migration.Migration{
		Name:         "00004_069662",
		Dependencies: []string{"00003_c93741"},
		Operations: []migration.Operation{
			migration.AlterTableOperation{
				Name: "user",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableDropColumn{Name: "secret_hash"},
					migration.AlterTableAddColumn{
						Name:     "secret",
						Type:     "text",
						Default:  "''",
						Nullable: false,
					},
				},
			},
		},
	})
}
