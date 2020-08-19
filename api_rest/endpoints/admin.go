package endpoints

import (
	"github.com/Eugenill/SmartScooter/api_rest/handlers"
	"github.com/Eugenill/SmartScooter/api_rest/handlers/user"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/mqtt_pub"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/mqtt_sub"
	"github.com/gin-gonic/gin"
)

func AddAdmin(group *gin.RouterGroup, mqttConf mqtt_sub.MQTTConfig) {
	group.GET("/publish_detection", mqtt_pub.PublishDetection(mqttConf, "detection_example"))
	group.GET("/vehicle", handlers.GetVehicles("Ford Mustang"))
	group.POST("/create_user", user.AdminCreateUser())
	group.POST("/edit_user", user.AdminEditUser())
	group.POST("/delete_user", user.AdminDeleteUser())
	group.POST("/get_users", user.AdminGetUsers())
}
