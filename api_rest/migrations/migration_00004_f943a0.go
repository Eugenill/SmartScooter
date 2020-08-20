package migrations

import "github.com/sqlbunny/sqlbunny/runtime/migration"

func init() {
	Store.Register(&migration.Migration{
		Name:         "00004_f943a0",
		Dependencies: []string{"00003_dea14c"},
		Operations: []migration.Operation{
			migration.AlterTableOperation{
				Name: "helmet",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableSetNull{Name: "last_ping"},
					migration.AlterTableDropDefault{Name: "last_ping"},
				},
			},
			migration.AlterTableOperation{
				Name: "ride_detection",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableSetNotNull{Name: "detection__traffic_sign"},
					migration.AlterTableSetDefault{Name: "detection__traffic_sign", Default: "0"},
				},
			},
		},
	})
}
