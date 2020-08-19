package migrations

import "github.com/sqlbunny/sqlbunny/runtime/migration"

func init() {
	Store.Register(&migration.Migration{
		Name:         "00002_ea776f",
		Dependencies: []string{"00001_4301da"},
		Operations: []migration.Operation{
			migration.AlterTableOperation{
				Name: "helmet",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableDropColumn{Name: "vehicle_zone"},
				},
			},
		},
	})
}
