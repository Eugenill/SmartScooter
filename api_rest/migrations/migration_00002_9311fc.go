package migrations

import "github.com/sqlbunny/sqlbunny/runtime/migration"

func init() {
	Store.Register(&migration.Migration{
		Name:         "00002_9311fc",
		Dependencies: []string{"00001_0f4f57"},
		Operations: []migration.Operation{
			migration.AlterTableOperation{
				Name: "all_detection",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableDropForeignKey{
						Name: "all_detection___ride_id___fkey",
					},
					migration.AlterTableDropForeignKey{
						Name: "all_detection___user_id___fkey",
					},
				},
			},
			migration.AlterTableOperation{
				Name: "all_detection",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableDropPrimaryKey{},
				},
			},
			migration.DropTableOperation{
				Name: "all_detection",
			},
			migration.CreateTableOperation{
				Name: "ride_detection",
				Columns: []migration.Column{
					migration.Column{Name: "id", Type: "bytea", Default: "'\\x000000000000'", Nullable: false},
					migration.Column{Name: "ride_id", Type: "bytea", Default: "'\\x000000000000'", Nullable: false},
					migration.Column{Name: "user_id", Type: "bytea", Default: "'\\x000000000000'", Nullable: false},
					migration.Column{Name: "detection__traffic_light", Type: "text", Default: "", Nullable: true},
					migration.Column{Name: "detection__obstacle", Type: "text", Default: "", Nullable: true},
					migration.Column{Name: "detection__traffic_sign", Type: "integer", Default: "", Nullable: true},
					migration.Column{Name: "detection__location__latitude", Type: "double precision", Default: "0", Nullable: false},
					migration.Column{Name: "detection__location__longitude", Type: "double precision", Default: "0", Nullable: false},
					migration.Column{Name: "detection__location__accuracy", Type: "double precision", Default: "0", Nullable: false},
					migration.Column{Name: "detection__detected_at", Type: "timestamptz", Default: "'0001-01-01 00:00:00+00'", Nullable: false},
				},
			},
			migration.AlterTableOperation{
				Name: "ride_detection",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableCreatePrimaryKey{
						Columns: []string{"id"},
					},
				},
			},
			migration.AlterTableOperation{
				Name: "ride_detection",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableCreateForeignKey{
						Name:           "ride_detection___ride_id___fkey",
						Columns:        []string{"ride_id"},
						ForeignTable:   "ride",
						ForeignColumns: []string{"id"},
					},
					migration.AlterTableCreateForeignKey{
						Name:           "ride_detection___user_id___fkey",
						Columns:        []string{"user_id"},
						ForeignTable:   "user",
						ForeignColumns: []string{"id"},
					},
				},
			},
		},
	})
}
