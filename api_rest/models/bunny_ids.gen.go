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
	case "pa":
		return PathIDFromString(s)
	case "h":
		return HelmetIDFromString(s)
	case "v":
		return VehicleIDFromString(s)
	case "r":
		return RideIDFromString(s)
	case "usr":
		return UserIDFromString(s)
	case "rd":
		return RideDetectionIDFromString(s)
	}
	return nil, fmt.Errorf("Unknown ID type %s", parts[0])
}

var idPrefixes = map[string]string{
	"pa":  "path_id",
	"h":   "helmet_id",
	"v":   "vehicle_id",
	"r":   "ride_id",
	"usr": "user_id",
	"rd":  "ride_detection_id",
}
