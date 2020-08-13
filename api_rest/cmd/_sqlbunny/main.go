package main

import (
	"github.com/Eugenill/SmartScooter/api_rest/migrations"
	bunnyid "github.com/sqlbunny/bunnyid/gen"
	. "github.com/sqlbunny/sqlbunny/gen/core"
	"github.com/sqlbunny/sqlbunny/gen/migration"
	"github.com/sqlbunny/sqlbunny/gen/stdtypes"
)

//USAGE:  go run ./cmd/_sqlbunny/main.go gen
// go run ./cmd/_sqlbunny/main.go migration gen

func main() {
	Run(
		&migration.Plugin{
			Store: &migrations.Store,
		},
		&stdtypes.Plugin{},
		&bunnyid.Plugin{},

		Type("point", Struct(
			Field("latitude", "float64"),
			Field("longitude", "float64"),
			Field("accuracy", "float64"),
		)),

		Type("vehicle_zone", Enum(
			"carril_bici",
			"carril_zona30",
			"acera",
		)),

		Type("helmet_status", Enum(
			"connected",
			"not_connected",
			"connection_error",
		)),
		Type("traffic_sign", Enum(
			"stop",
			"yield",
			"vel_10",
			"vel_20",
			"vel_30",
		)),

		Type("detection", Struct(
			Field("traffic_light", "string", Null),
			Field("obstacle", "string", Null),
			Field("traffic_sign", "traffic_sign", Null),
			Field("location", "point"),
			Field("detected_at", "time"),
		)),

		Type("helmet_id", bunnyid.ID{Prefix: "h"}),
		Model("helmet",
			Field("id", "helmet_id", PrimaryKey),
			Field("vehicle_zone", "vehicle_zone"),
			Field("last_ping", "time"),
			Field("helmet_status", "helmet_status"),
		),

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
		Type("ride_detection_id", bunnyid.ID{Prefix: "rd"}),
		Model("ride_detection",
			Field("id", "ride_detection_id", PrimaryKey),
			Field("ride_id", "ride_id", ForeignKey("ride")),
			Field("user_id", "user_id", ForeignKey("user")),
			Field("detection", "detection"),
		),
	)
}
