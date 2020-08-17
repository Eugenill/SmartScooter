package migrations

import "github.com/sqlbunny/sqlbunny/runtime/migration"

func init() {
	Store.Register(&migration.Migration{
		Name:         "00003_c93741",
		Dependencies: []string{"00002_05365d"},
		Operations: []migration.Operation{
			migration.AlterTableOperation{
				Name: "user",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableDropUnique{
						Name: "user___login___key",
					},
				},
			},
			migration.DropIndexOperation{
				Name:      "user",
				IndexName: "user___login___idx",
			},
			migration.AlterTableOperation{
				Name: "user",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableDropColumn{Name: "login"},
					migration.AlterTableAddColumn{
						Name:     "username",
						Type:     "text",
						Default:  "''",
						Nullable: false,
					},
					migration.AlterTableAddColumn{
						Name:     "admin",
						Type:     "boolean",
						Default:  "false",
						Nullable: false,
					},
				},
			},
			migration.CreateIndexOperation{
				Name:      "user",
				IndexName: "user___username___idx",
				Columns:   []string{"username"},
			},
			migration.AlterTableOperation{
				Name: "user",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableCreateUnique{
						Name:    "user___username___key",
						Columns: []string{"username"},
					},
				},
			},
		},
	})
}
