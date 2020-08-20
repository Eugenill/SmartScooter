package helmet

import (
	"context"
	"github.com/Eugenill/SmartScooter/api_rest/models"
)

//this function context must contain the DB
func CheckHelmet(ctx context.Context, id models.HelmetID) (bool, error) {
	_, err := models.FindHelmet(ctx, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

type CreateHelmet struct {
	Name string `json:"name"`
}
