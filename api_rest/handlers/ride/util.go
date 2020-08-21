package ride

import (
	"context"
	"github.com/Eugenill/SmartScooter/api_rest/models"
	"github.com/sqlbunny/geo"
	"github.com/sqlbunny/sqlbunny/runtime/bunny"
	"github.com/sqlbunny/sqlbunny/runtime/qm"
)

func CheckRide(ctx context.Context, id models.RideID) (bool, error) {
	_, err := models.FindRide(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func CalcPath(ctx context.Context, r *models.Ride) (geo.LineStringM, float32, error) {
	rideDetections, err := models.RideDetections(
		qm.Where("ride_id = ?", r.ID),
	).All(ctx)
	if err != nil {
		return geo.LineStringM{}, 0, err
	}

	res := geo.LineStringM{}
	for _, rideDetection := range rideDetections {
		if rideDetection.Detection.Location.X != 0 && rideDetection.Detection.Location.Y != 0 {
			res = append(res, geo.PointM{
				X: rideDetection.Detection.Location.X,
				Y: rideDetection.Detection.Location.Y,
				M: rideDetection.CreatedAt.Sub(r.StartedAt).Seconds(),
			})
		}
	}
	distance, err := CalcDistance(ctx, res)
	if err != nil {
		return geo.LineStringM{}, 0, err
	}
	return res, distance, nil
}

func CalcDistance(ctx context.Context, path geo.LineStringM) (float32, error) {
	var res float32
	err := bunny.QueryRow(ctx, "SELECT ST_Length($1::geography(LineStringM, 4326))", path).Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}
