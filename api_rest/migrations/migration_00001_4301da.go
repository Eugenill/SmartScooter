package migrations

import "github.com/sqlbunny/sqlbunny/runtime/migration"

func init() {
	Store.Register(&migration.Migration{
		Name:         "00001_4301da",
		Dependencies: []string{},
		Operations: []migration.Operation{
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
					migration.Column{Name: "detection__detection_zone", Type: "integer", Default: "0", Nullable: false},
				},
			},
			migration.CreateTableOperation{
				Name: "path",
				Columns: []migration.Column{
					migration.Column{Name: "id", Type: "bytea", Default: "'\\x000000000000'", Nullable: false},
					migration.Column{Name: "ride_id", Type: "bytea", Default: "'\\x000000000000'", Nullable: false},
					migration.Column{Name: "point__latitude", Type: "double precision", Default: "0", Nullable: false},
					migration.Column{Name: "point__longitude", Type: "double precision", Default: "0", Nullable: false},
					migration.Column{Name: "point__accuracy", Type: "double precision", Default: "0", Nullable: false},
				},
			},
			migration.CreateTableOperation{
				Name: "user",
				Columns: []migration.Column{
					migration.Column{Name: "id", Type: "bytea", Default: "'\\x000000000000'", Nullable: false},
					migration.Column{Name: "username", Type: "text", Default: "''", Nullable: false},
					migration.Column{Name: "secret", Type: "text", Default: "''", Nullable: false},
					migration.Column{Name: "contact_email", Type: "text", Default: "''", Nullable: false},
					migration.Column{Name: "admin", Type: "boolean", Default: "false", Nullable: false},
					migration.Column{Name: "phone_number", Type: "text", Default: "''", Nullable: false},
					migration.Column{Name: "created_at", Type: "timestamptz", Default: "'0001-01-01 00:00:00+00'", Nullable: false},
					migration.Column{Name: "is_deleted", Type: "boolean", Default: "false", Nullable: false},
					migration.Column{Name: "deleted_at", Type: "timestamptz", Default: "", Nullable: true},
				},
			},
			migration.CreateTableOperation{
				Name: "ride",
				Columns: []migration.Column{
					migration.Column{Name: "id", Type: "bytea", Default: "'\\x000000000000'", Nullable: false},
					migration.Column{Name: "vehicle_id", Type: "bytea", Default: "'\\x000000000000'", Nullable: false},
					migration.Column{Name: "user_id", Type: "bytea", Default: "'\\x000000000000'", Nullable: false},
					migration.Column{Name: "path_id", Type: "bytea", Default: "'\\x000000000000'", Nullable: false},
					migration.Column{Name: "distance", Type: "real", Default: "0", Nullable: false},
					migration.Column{Name: "duration", Type: "integer", Default: "0", Nullable: false},
					migration.Column{Name: "started_at", Type: "timestamptz", Default: "'0001-01-01 00:00:00+00'", Nullable: false},
					migration.Column{Name: "finished_at", Type: "timestamptz", Default: "", Nullable: true},
				},
			},
			migration.CreateTableOperation{
				Name: "vehicle",
				Columns: []migration.Column{
					migration.Column{Name: "id", Type: "bytea", Default: "'\\x000000000000'", Nullable: false},
					migration.Column{Name: "current_ride_id", Type: "bytea", Default: "'\\x000000000000'", Nullable: false},
					migration.Column{Name: "last_ride_id", Type: "bytea", Default: "'\\x000000000000'", Nullable: false},
					migration.Column{Name: "current_user_id", Type: "bytea", Default: "'\\x000000000000'", Nullable: false},
					migration.Column{Name: "last_user_id", Type: "bytea", Default: "'\\x000000000000'", Nullable: false},
					migration.Column{Name: "number_plate", Type: "text", Default: "''", Nullable: false},
					migration.Column{Name: "helmet_id", Type: "bytea", Default: "'\\x000000000000'", Nullable: false},
				},
			},
			migration.CreateIndexOperation{
				Name:      "user",
				IndexName: "user___username___idx",
				Columns:   []string{"username"},
			},
			migration.CreateIndexOperation{
				Name:      "user",
				IndexName: "user___contact_email___idx",
				Columns:   []string{"contact_email"},
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
				Name: "ride_detection",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableCreatePrimaryKey{
						Columns: []string{"id"},
					},
				},
			},
			migration.AlterTableOperation{
				Name: "path",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableCreatePrimaryKey{
						Columns: []string{"id"},
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
						Name:    "user___username___key",
						Columns: []string{"username"},
					},
				},
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
			migration.AlterTableOperation{
				Name: "path",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableCreateForeignKey{
						Name:           "path___ride_id___fkey",
						Columns:        []string{"ride_id"},
						ForeignTable:   "ride",
						ForeignColumns: []string{"id"},
					},
				},
			},
			migration.AlterTableOperation{
				Name: "ride",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableCreateForeignKey{
						Name:           "ride___vehicle_id___fkey",
						Columns:        []string{"vehicle_id"},
						ForeignTable:   "vehicle",
						ForeignColumns: []string{"id"},
					},
					migration.AlterTableCreateForeignKey{
						Name:           "ride___user_id___fkey",
						Columns:        []string{"user_id"},
						ForeignTable:   "user",
						ForeignColumns: []string{"id"},
					},
					migration.AlterTableCreateForeignKey{
						Name:           "ride___path_id___fkey",
						Columns:        []string{"path_id"},
						ForeignTable:   "path",
						ForeignColumns: []string{"id"},
					},
				},
			},
			migration.AlterTableOperation{
				Name: "vehicle",
				Ops: []migration.AlterTableSuboperation{
					migration.AlterTableCreateForeignKey{
						Name:           "vehicle___current_ride_id___fkey",
						Columns:        []string{"current_ride_id"},
						ForeignTable:   "ride",
						ForeignColumns: []string{"id"},
					},
					migration.AlterTableCreateForeignKey{
						Name:           "vehicle___last_ride_id___fkey",
						Columns:        []string{"last_ride_id"},
						ForeignTable:   "ride",
						ForeignColumns: []string{"id"},
					},
					migration.AlterTableCreateForeignKey{
						Name:           "vehicle___current_user_id___fkey",
						Columns:        []string{"current_user_id"},
						ForeignTable:   "user",
						ForeignColumns: []string{"id"},
					},
					migration.AlterTableCreateForeignKey{
						Name:           "vehicle___last_user_id___fkey",
						Columns:        []string{"last_user_id"},
						ForeignTable:   "user",
						ForeignColumns: []string{"id"},
					},
					migration.AlterTableCreateForeignKey{
						Name:           "vehicle___helmet_id___fkey",
						Columns:        []string{"helmet_id"},
						ForeignTable:   "helmet",
						ForeignColumns: []string{"id"},
					},
				},
			},
		},
	})
}
