package migrations

import "github.com/sqlbunny/sqlbunny/runtime/migration"

func init() {
	Store.Register(&migration.Migration{
		Name:         "00007_f23ba3",
		Dependencies: []string{"00006_77cf2e"},
		Operations: []migration.Operation{
			migration.DropIndexOperation{
				Name:      "vehicle",
				IndexName: "vehicle___helmet_id___idx",
			},
			migration.AlterTableOperation{
				Name: "vehicle",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableAddColumn{
						Name:     "iot_device_id",
						Type:     "bytea",
						Default:  "'\\x000000000000'",
						Nullable: false,
					},
				},
			},
			migration.CreateTableOperation{
				Name: "iot_device",
				Columns: []migration.Column{
					migration.Column{Name: "id", Type: "bytea", Default: "'\\x000000000000'", Nullable: false},
					migration.Column{Name: "name", Type: "text", Default: "''", Nullable: false},
					migration.Column{Name: "last_ping", Type: "timestamptz", Default: "", Nullable: true},
					migration.Column{Name: "iot_device_status", Type: "integer", Default: "0", Nullable: false},
				},
			},
			migration.AlterTableOperation{
				Name: "iot_device",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableCreatePrimaryKey{
						Columns: []string{"id"},
					},
					migration.AlterTableCreateUnique{
						Name:    "iot_device___name___key",
						Columns: []string{"name"},
					},
				},
			},
			migration.AlterTableOperation{
				Name: "vehicle",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableCreateForeignKey{
						Name:           "vehicle___iot_device_id___fkey",
						Columns:        []string{"iot_device_id"},
						ForeignTable:   "iot_device",
						ForeignColumns: []string{"id"},
					},
				},
			},
		},
	})
}
