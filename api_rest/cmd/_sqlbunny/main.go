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

		Type("point", BaseType{
			Go:     "github.com/sqlbunny/geo.Point",
			GoNull: "github.com/sqlbunny/geo/null.Point",
			Postgres: SQLType{
				Type: "geography(Point, 4326)",
			},
		}),

		Type("line_string_m", BaseType{
			Go:     "github.com/sqlbunny/geo.LineStringM",
			GoNull: "github.com/sqlbunny/geo/null.LineStringM",
			Postgres: SQLType{
				Type: "geography(LineStringM, 4326)",
			},
		}),

		Type("vehicle_zone", Enum(
			"None",
			"carril_bici",
			"carril_zona30",
			"acera",
		)),

		Type("device_status", Enum(
			"None",
			"connected",
			"not_connected",
			"connection_error",
		)),
		Type("traffic_sign", Enum(
			"None",
			"stop",
			"yield",
			"vel_10",
			"vel_20",
			"vel_30",
			"people",
		)),

		Type("detection", Struct(
			Field("traffic_light", "string", Null),
			Field("obstacle", "string", Null),
			Field("traffic_sign", "traffic_sign"),
			Field("location", "point"),
			Field("detected_at", "time"),
			Field("detection_zone", "vehicle_zone"),
		)),

		Type("helmet_id", bunnyid.ID{Prefix: "h"}),
		Model("helmet",
			Field("id", "helmet_id", PrimaryKey),
			Field("name", "string", Unique),
			Field("last_ping", "time", Null),
			Field("helmet_status", "device_status"),
		),

		Type("iot_device_id", bunnyid.ID{Prefix: "iot"}),
		Model("iot_device",
			Field("id", "iot_device_id", PrimaryKey),
			Field("name", "string", Unique),
			Field("last_ping", "time", Null),
			Field("iot_device_status", "device_status"),
		),

		Type("vehicle_id", bunnyid.ID{Prefix: "v"}),
		Model("vehicle",
			Field("id", "vehicle_id", PrimaryKey),
			Field("current_ride_id", "ride_id", Null, ForeignKey("ride")),
			Field("last_ride_id", "ride_id", Null, ForeignKey("ride")),
			Field("current_user_id", "user_id", Null, ForeignKey("user")),
			Field("last_user_id", "user_id", Null, ForeignKey("user")),
			Field("number_plate", "string", Unique),
			Field("helmet_id", "helmet_id", ForeignKey("helmet")),
			Field("iot_device_id", "iot_device_id", ForeignKey("iot_device")),
		),

		Type("ride_id", bunnyid.ID{Prefix: "r"}),
		Model("ride",
			Field("id", "ride_id", PrimaryKey),
			Field("vehicle_id", "vehicle_id", ForeignKey("vehicle")),
			Field("user_id", "user_id", ForeignKey("user")),
			Field("path", "line_string_m", Null),
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
			Field("contact_email", "string", Index, Unique),
			Field("admin", "bool"),
			Field("phone_number", "string", Unique),
			Field("created_at", "time"),
			Field("is_deleted", "bool"),
			Field("deleted_at", "time", Null),
		),
		Type("ride_detection_id", bunnyid.ID{Prefix: "rd"}),
		Model("ride_detection",
			Field("id", "ride_detection_id", PrimaryKey),
			Field("ride_id", "ride_id", ForeignKey("ride")),
			Field("detection", "detection"),
			Field("created_at", "time"),
		),
	)
}
