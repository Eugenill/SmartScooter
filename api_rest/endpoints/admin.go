package endpoints

import (
	"github.com/Eugenill/SmartScooter/api_rest/handlers/helmet"
	"github.com/Eugenill/SmartScooter/api_rest/handlers/iot_dev"
	"github.com/Eugenill/SmartScooter/api_rest/handlers/ride"
	"github.com/Eugenill/SmartScooter/api_rest/handlers/user"
	"github.com/Eugenill/SmartScooter/api_rest/handlers/vehicle"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/mqtt_pub"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/mqtt_sub"
	"github.com/gin-gonic/gin"
)

func AddAdmin(group *gin.RouterGroup, mqttConf mqtt_sub.MQTTConfig) {
	group.GET("/publish_detection", mqtt_pub.PublishDetection(mqttConf, "detection_example"))

	usr := group.Group("/user")
	usr.POST("/create", user.AdminCreateUser())
	usr.POST("/edit", user.AdminEditUser())
	usr.POST("/delete", user.AdminDeleteUser())
	usr.POST("/get", user.AdminGetUsers())

	r := group.Group("/ride")
	r.GET("/create", ride.CreateRide())
	r.GET("/finish", ride.FinishRide())
	r.POST("/get", ride.AdminGetRides())

	veh := group.Group("/vehicle")
	veh.POST("/create", vehicle.CreateVehicle())
	veh.POST("/delete", vehicle.AdminDeleteVehicle())
	veh.POST("/get", vehicle.AdminGetVehicles())

	helm := group.Group("/helmet")
	helm.POST("/create", helmet.CreateHelmet())
	helm.POST("/delete", helmet.AdminDeleteHelmet())
	helm.POST("/get", helmet.AdminGetHelmets())

	iot := group.Group("/iot_dev")
	iot.POST("/create", iot_dev.CreateIotDev())
	iot.POST("/delete", iot_dev.AdminDeleteIotDev())
	iot.POST("/get", iot_dev.AdminGetIotDevs())

	//	rd := group.Group("/ride_detections")
	//rd.GET("/get", ride.GetDetections())

}
