package main

import (
	bunnyid "github.com/sqlbunny/bunnyid/gen"
	. "github.com/sqlbunny/sqlbunny/gen/core"
	"github.com/sqlbunny/sqlbunny/gen/migration"
	"github.com/sqlbunny/sqlbunny/gen/stdtypes"
)

func main() {
	Run(
		&stdtypes.Plugin{},
		&migration.Plugin{},
		&bunnyid.Plugin{},

		Type("point", BaseType{
			Go:     "github.com/sqlbunny/geo.Point",
			GoNull: "github.com/sqlbunny/geo/null.Point",
			Postgres: SQLType{
				Type: "geography(Point, 4326)",
			},
		}),
		Type("vehicle_zone", Enum(
			"carril_bici",
			"carril_zona30",
			"acera",
		)),

		Type("vehicle_id", bunnyid.ID{Prefix: "v"}),
		Model("vehicle",
			Field("id", "string", PrimaryKey),
			Field("token", "string"),
			Field("number_plate", "string"),
			Field("vehicle_zone", "vehicle_zone"),
			Field("helmet_id", "helmet_id", ForeignKey("helmet"), Index),
		),

		Type("ride_id", bunnyid.ID{Prefix: "r"}),
		Model("ride",
			Field("id", "ride_id", PrimaryKey),
			Field("distance", "float32"),
			Field("started_at", "time", Index),
			Field("finished_at", "time", Null),
		),

		Type("user_id", bunnyid.ID{Prefix: "usr"}),
		Model("user",
			Field("id", "user_id", PrimaryKey),
			Field("login", "string", Unique),
			Field("secret_hash", "string"),
			Field("contact_email", "string"),
			Field("is_deleted", "bool"),
			Field("deleted_at", "time", Null),
		),

		Type("helmet_status", Enum(
			"conected",
			"not_connected",
			"connection_error",
		)),

		Type("detection", Struct(
			Field("traffic_light", "string", Null),
			Field("obstacle", "string", Null),
			Field("location", "point"),
			Field("detected_at", "time"),
		)),

		Model("all_detection",
			Field("ride_id", "ride_id", PrimaryKey, ForeignKey("ride")),
			Field("user_id", "user_id", ForeignKey("user")),
			Field("detection", "detection"),
		),

		Type("helmet_id", bunnyid.ID{Prefix: "h"}),
		Model("helmet",
			Field("id", "helmet_id", PrimaryKey),
			Field("vehicle_zone", "vehicle_zone"),
			Field("last_ping", "time"),
			Field("helmet_status", "helmet_status"),
		),
	)
}
