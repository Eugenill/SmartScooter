package migrations

import "github.com/sqlbunny/sqlbunny/runtime/migration"

func init() {
	Store.Register(&migration.Migration{
		Name:         "00006_77cf2e",
		Dependencies: []string{"00005_0479af"},
		Operations: []migration.Operation{
			migration.AlterTableOperation{
				Name: "user",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableCreateUnique{
						Name:    "user___contact_email___key",
						Columns: []string{"contact_email"},
					},
					migration.AlterTableCreateUnique{
						Name:    "user___phone_number___key",
						Columns: []string{"phone_number"},
					},
				},
			},
			migration.AlterTableOperation{
				Name: "vehicle",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableCreateUnique{
						Name:    "vehicle___number_plate___key",
						Columns: []string{"number_plate"},
					},
				},
			},
			migration.AlterTableOperation{
				Name: "helmet",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableCreateUnique{
						Name:    "helmet___name___key",
						Columns: []string{"name"},
					},
				},
			},
		},
	})
}
