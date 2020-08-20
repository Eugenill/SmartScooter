package migrations

import "github.com/sqlbunny/sqlbunny/runtime/migration"

func init() {
	Store.Register(&migration.Migration{
		Name:         "00003_dea14c",
		Dependencies: []string{"00002_ea776f"},
		Operations: []migration.Operation{
			migration.AlterTableOperation{
				Name: "vehicle",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableSetNull{Name: "current_ride_id"},
					migration.AlterTableDropDefault{Name: "current_ride_id"},
					migration.AlterTableSetNull{Name: "last_ride_id"},
					migration.AlterTableDropDefault{Name: "last_ride_id"},
					migration.AlterTableSetNull{Name: "current_user_id"},
					migration.AlterTableDropDefault{Name: "current_user_id"},
					migration.AlterTableSetNull{Name: "last_user_id"},
					migration.AlterTableDropDefault{Name: "last_user_id"},
				},
			},
		},
	})
}
