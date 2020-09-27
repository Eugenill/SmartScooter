package reby_scooter

import (
	"github.com/Eugenill/SmartScooter/api_rest/handlers/reby"
	"github.com/Eugenill/SmartScooter/api_rest/models"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RideCall(ctx *gin.Context, vID models.VehicleID, meta interface{}) (ginErr *gin.Error) {
	req, err := http.NewRequest("POST", reby.RebyHost+reby.RebyRide, nil)
	if err != nil {
		_, ginErr = errors.New(ctx, "ride call request creation failed", gin.ErrorTypePrivate, meta)
		return ginErr
	} else {
		req.Header = reby.SetHeaders(vID.String(), reby.BearerETSEIB)
	}
	client := http.DefaultClient
	_, err = client.Do(req)

	if err != nil {
		_, ginErr = errors.New(ctx, "ride call to the scooter failed", gin.ErrorTypePrivate, meta)
		return ginErr
	}
	return nil
}

func UnlockCall(ctx *gin.Context, vID models.VehicleID, meta interface{}) (ginErr *gin.Error) {
	req, err := http.NewRequest("POST", reby.RebyHost+reby.RebyUnlock, nil)
	if err != nil {
		_, ginErr = errors.New(ctx, "unlock call request creation failed", gin.ErrorTypePrivate, meta)
		return ginErr
	} else {
		req.Header = reby.SetHeaders(vID.String(), reby.BearerETSEIB)
	}
	client := http.DefaultClient
	_, err = client.Do(req)

	if err != nil {
		_, ginErr = errors.New(ctx, "unlock call to the scooter failed", gin.ErrorTypePrivate, meta)
		return ginErr
	}
	return nil
}

func LockCall(ctx *gin.Context, vID models.VehicleID, meta interface{}) (ginErr *gin.Error) {
	req, err := http.NewRequest("POST", reby.RebyHost+reby.RebyLock, nil)
	if err != nil {
		_, ginErr = errors.New(ctx, "lock call request creation failed", gin.ErrorTypePrivate, meta)
		return ginErr
	} else {
		req.Header = reby.SetHeaders(vID.String(), reby.BearerETSEIB)
	}
	client := http.DefaultClient
	_, err = client.Do(req)

	if err != nil {
		_, ginErr = errors.New(ctx, "lock call to the scooter failed", gin.ErrorTypePrivate, meta)
		return ginErr
	}
	return nil
}

func SoundCall(ctx *gin.Context, vID models.VehicleID, meta interface{}) (ginErr *gin.Error) {
	req, err := http.NewRequest("POST", reby.RebyHost+reby.RebySound, nil)
	if err != nil {
		_, ginErr = errors.New(ctx, "sound call request creation failed", gin.ErrorTypePrivate, meta)
		return ginErr
	} else {
		req.Header = reby.SetHeaders(vID.String(), reby.BearerETSEIB)
	}
	client := http.DefaultClient
	_, err = client.Do(req)

	if err != nil {
		_, ginErr = errors.New(ctx, "sound call to the scooter failed", gin.ErrorTypePrivate, meta)
		return ginErr
	}
	return nil
}

func StatusCall(ctx *gin.Context, vID models.VehicleID, meta interface{}) (ginErr *gin.Error, response *http.Response) {
	req, err := http.NewRequest("POST", reby.RebyHost+reby.RebyStatus, nil)
	if err != nil {
		_, ginErr = errors.New(ctx, "status call request creation failed", gin.ErrorTypePrivate, meta)
		return ginErr, nil
	} else {
		req.Header = reby.SetHeaders(vID.String(), reby.BearerETSEIB)
	}
	client := http.DefaultClient
	response, err = client.Do(req)

	if err != nil {
		_, ginErr = errors.New(ctx, "status call to the scooter failed", gin.ErrorTypePrivate, meta)
		return ginErr, nil
	}
	//defer response.Body.Close() //to put in the function that calls it
	return nil, response
}
