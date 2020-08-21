package models

import (
	"fmt"
	"strings"
)

func IDFromString(s string) (interface{}, error) {
	parts := strings.Split(s, "_")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Wrong parts count, expected 2 got %d", len(parts))
	}
	switch parts[0] {
	case "h":
		return HelmetIDFromString(s)
	case "r":
		return RideIDFromString(s)
	case "usr":
		return UserIDFromString(s)
	case "rd":
		return RideDetectionIDFromString(s)
	case "iot":
		return IotDeviceIDFromString(s)
	case "v":
		return VehicleIDFromString(s)
	}
	return nil, fmt.Errorf("Unknown ID type %s", parts[0])
}

var idPrefixes = map[string]string{
	"h":   "helmet_id",
	"r":   "ride_id",
	"usr": "user_id",
	"rd":  "ride_detection_id",
	"iot": "iot_device_id",
	"v":   "vehicle_id",
}
