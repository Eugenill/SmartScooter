package migrations

import "github.com/sqlbunny/sqlbunny/runtime/migration"

func init() {
	Store.Register(&migration.Migration{
		Name:         "00001_0f4f57",
		Dependencies: []string{},
		Operations: []migration.Operation{
			migration.CreateTableOperation{
				Name: "user",
				Columns: []migration.Column{
					migration.Column{Name: "id", Type: "bytea", Default: "'\\x000000000000'", Nullable: false},
					migration.Column{Name: "login", Type: "text", Default: "''", Nullable: false},
					migration.Column{Name: "secret_hash", Type: "text", Default: "''", Nullable: false},
					migration.Column{Name: "contact_email", Type: "text", Default: "''", Nullable: false},
					migration.Column{Name: "is_deleted", Type: "boolean", Default: "false", Nullable: false},
					migration.Column{Name: "deleted_at", Type: "timestamptz", Default: "", Nullable: true},
				},
			},
			migration.CreateTableOperation{
				Name: "ride",
				Columns: []migration.Column{
					migration.Column{Name: "id", Type: "bytea", Default: "'\\x000000000000'", Nullable: false},
					migration.Column{Name: "distance", Type: "real", Default: "0", Nullable: false},
					migration.Column{Name: "started_at", Type: "timestamptz", Default: "'0001-01-01 00:00:00+00'", Nullable: false},
					migration.Column{Name: "finished_at", Type: "timestamptz", Default: "", Nullable: true},
				},
			},
			migration.CreateTableOperation{
				Name: "vehicle",
				Columns: []migration.Column{
					migration.Column{Name: "id", Type: "text", Default: "''", Nullable: false},
					migration.Column{Name: "token", Type: "text", Default: "''", Nullable: false},
					migration.Column{Name: "number_plate", Type: "text", Default: "''", Nullable: false},
					migration.Column{Name: "vehicle_zone", Type: "integer", Default: "0", Nullable: false},
					migration.Column{Name: "helmet_id", Type: "bytea", Default: "'\\x000000000000'", Nullable: false},
				},
			},
			migration.CreateTableOperation{
				Name: "helmet",
				Columns: []migration.Column{
					migration.Column{Name: "id", Type: "bytea", Default: "'\\x000000000000'", Nullable: false},
					migration.Column{Name: "vehicle_zone", Type: "integer", Default: "0", Nullable: false},
					migration.Column{Name: "last_ping", Type: "timestamptz", Default: "'0001-01-01 00:00:00+00'", Nullable: false},
					migration.Column{Name: "helmet_status", Type: "integer", Default: "0", Nullable: false},
				},
			},
			migration.CreateTableOperation{
				Name: "all_detection",
				Columns: []migration.Column{
					migration.Column{Name: "ride_id", Type: "bytea", Default: "'\\x000000000000'", Nullable: false},
					migration.Column{Name: "user_id", Type: "bytea", Default: "'\\x000000000000'", Nullable: false},
					migration.Column{Name: "detection__traffic_light", Type: "text", Default: "", Nullable: true},
					migration.Column{Name: "detection__obstacle", Type: "text", Default: "", Nullable: true},
					migration.Column{Name: "detection__location__latitude", Type: "double precision", Default: "0", Nullable: false},
					migration.Column{Name: "detection__location__longitude", Type: "double precision", Default: "0", Nullable: false},
					migration.Column{Name: "detection__location__accuracy", Type: "double precision", Default: "0", Nullable: false},
					migration.Column{Name: "detection__detected_at", Type: "timestamptz", Default: "'0001-01-01 00:00:00+00'", Nullable: false},
				},
			},
			migration.CreateIndexOperation{
				Name:      "ride",
				IndexName: "ride___started_at___idx",
				Columns:   []string{"started_at"},
			},
			migration.CreateIndexOperation{
				Name:      "vehicle",
				IndexName: "vehicle___helmet_id___idx",
				Columns:   []string{"helmet_id"},
			},
			migration.AlterTableOperation{
				Name: "ride",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableCreatePrimaryKey{
						Columns: []string{"id"},
					},
				},
			},
			migration.AlterTableOperation{
				Name: "vehicle",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableCreatePrimaryKey{
						Columns: []string{"id"},
					},
				},
			},
			migration.AlterTableOperation{
				Name: "helmet",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableCreatePrimaryKey{
						Columns: []string{"id"},
					},
				},
			},
			migration.AlterTableOperation{
				Name: "all_detection",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableCreatePrimaryKey{
						Columns: []string{"ride_id"},
					},
				},
			},
			migration.AlterTableOperation{
				Name: "user",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableCreatePrimaryKey{
						Columns: []string{"id"},
					},
					migration.AlterTableCreateUnique{
						Name:    "user___login___key",
						Columns: []string{"login"},
					},
				},
			},
			migration.AlterTableOperation{
				Name: "vehicle",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableCreateForeignKey{
						Name:           "vehicle___helmet_id___fkey",
						Columns:        []string{"helmet_id"},
						ForeignTable:   "helmet",
						ForeignColumns: []string{"id"},
					},
				},
			},
			migration.AlterTableOperation{
				Name: "all_detection",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableCreateForeignKey{
						Name:           "all_detection___ride_id___fkey",
						Columns:        []string{"ride_id"},
						ForeignTable:   "ride",
						ForeignColumns: []string{"id"},
					},
					migration.AlterTableCreateForeignKey{
						Name:           "all_detection___user_id___fkey",
						Columns:        []string{"user_id"},
						ForeignTable:   "user",
						ForeignColumns: []string{"id"},
					},
				},
			},
		},
	})
}
