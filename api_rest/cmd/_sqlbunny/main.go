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
			"people",
		)),
		Type("path_id", bunnyid.ID{Prefix: "pa"}),
		Model("path",
			Field("id", "path_id", PrimaryKey),
			Field("ride_id", "ride_id", ForeignKey("ride")),
			Field("point", "point"),
		),
		Type("detection", Struct(
			Field("traffic_light", "string", Null),
			Field("obstacle", "string", Null),
			Field("traffic_sign", "traffic_sign", Null),
			Field("location", "point"),
			Field("detected_at", "time"),
			Field("detection_zone", "vehicle_zone"),
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
			Field("id", "vehicle_id", PrimaryKey),
			Field("current_ride_id", "ride_id", ForeignKey("ride")),
			Field("last_ride_id", "ride_id", ForeignKey("ride")),
			Field("current_user_id", "user_id", ForeignKey("user")),
			Field("last_user_id", "user_id", ForeignKey("user")),
			Field("number_plate", "string"),
			Field("helmet_id", "helmet_id", ForeignKey("helmet"), Index),
		),

		Type("ride_id", bunnyid.ID{Prefix: "r"}),
		Model("ride",
			Field("id", "ride_id", PrimaryKey),
			Field("vehicle_id", "vehicle_id", ForeignKey("vehicle")),
			Field("user_id", "user_id", ForeignKey("user")),
			Field("path_id", "path_id", ForeignKey("path")),
			Field("distance", "float32"),
			Field("duration", "int32"),
			Field("started_at", "time", Index),
			Field("finished_at", "time", Null),
		),

		Type("user_id", bunnyid.ID{Prefix: "usr"}),
		Model("user",
			Field("id", "user_id", PrimaryKey),
			Field("username", "string", Unique, Index),
			Field("secret", "string"),
			Field("contact_email", "string", Index),
			Field("admin", "bool"),
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
