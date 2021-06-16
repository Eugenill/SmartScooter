package kernel

import (
	"github.com/Eugenill/SmartScooter/Zone_Classifier/backend/pkg/kernel/http/api"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// Service defines the methods that any service must comply to be managed by the kernel.
type Service interface {
	RegisterHTTPHandlers() []api.HTTPHandlerDefinition
	Init()
}

type Kernel struct {
	Engine      *gin.Engine
	Middlewares []gin.HandlerFunc
	Services    []Service
}

func New(engine *gin.Engine, middlewares []gin.HandlerFunc) *Kernel {
	return &Kernel{
		Engine:      engine,
		Middlewares: middlewares,
		Services:    make([]Service, 0),
	}
}

// Init will initialize the kernel with all those functions that are not suitable to be initialized or started on New()
func (k *Kernel) Init() error {
	k.initServices()
	k.registerHTTPEndpoints()
	return k.RunEngine()

}
func (k *Kernel) AddService(service Service) {
	k.Services = append(k.Services, service)
}

func (k *Kernel) registerHTTPEndpoints() {
	k.UseMiddlewares()
	for _, service := range k.Services {
		for _, def := range service.RegisterHTTPHandlers() {
			switch strings.ToUpper(def.Method) {
			case http.MethodGet:
				def.RouterGroup.GET(def.Endpoint, def.Handler)
			case http.MethodPost:
				def.RouterGroup.POST(def.Endpoint, def.Handler)
			case http.MethodPatch:
				def.RouterGroup.PATCH(def.Endpoint, def.Handler)
			case http.MethodPut:
				def.RouterGroup.PUT(def.Endpoint, def.Handler)
			case http.MethodDelete:
				def.RouterGroup.DELETE(def.Endpoint, def.Handler)
			default:
				panic("wrong http method " + def.Method)
			}
		}
	}
}

func (k *Kernel) UseMiddlewares() {
	for _, middleware := range k.Middlewares {
		k.Engine.Use(middleware)
	}
}
func (k *Kernel) initServices() {
	for _, service := range k.Services {
		service.Init()
	}
}

func (k *Kernel) RunEngine() error {
	err := k.Engine.Run("localhost:1234")
	return err
}
