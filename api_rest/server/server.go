package server

import (
	"encoding/base64"
	"fmt"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/mqtt_sub"
	"github.com/gin-gonic/gin"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
)

// RequestIDKey is the key that holds the unique request ID in a request context.
const RequestIDKey int = 0

var reqid uint64

type Router struct {
	Engine      *gin.Engine
	Middlewares []gin.HandlerFunc
	MqttConfig  mqtt_sub.MQTTConfig
	AdminRoute  *gin.RouterGroup
	PublicRoute *gin.RouterGroup
}

func (s *Router) AddMiddlewares() {
	for _, middleware := range s.Middlewares {
		s.Engine.Use(middleware)
	}
}

func (s *Router) RunServer() error {
	err := s.Engine.Run("localhost:1234")
	return err
}

// RequestID is a middleware that injects a request ID into the context of each
// request. A request ID is a string of the form "host.example.com/random-0001",
// where "random" is a base62 random string that uniquely identifies this go
// process, and where the last number is an atomically incremented request
// counter.
func RequestID() gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		r := ctx.Request
		requestID := r.Header.Get("X-Request-Id")
		if requestID == "" {
			myid := atomic.AddUint64(&reqid, 1)
			requestID = fmt.Sprintf("%s-%06d", prefixGen(), myid)
		}
		ctx.Set(strconv.Itoa(RequestIDKey), requestID)
		ctx.Next()
	}
	return fn
}

func prefixGen() string {
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}
	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		rand.Read(buf[:])
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}
	return fmt.Sprintf("%s/%s", hostname, b64[0:10])
}

/*
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		if params.IsOutputColor() {
			statusColor = params.StatusCodeColor()
			methodColor = params.MethodColor()
			resetColor = params.ResetColor()
		}

		if params.Latency > time.Minute {
			// Truncate in a golang < 1.8 safe way
			params.Latency = params.Latency - params.Latency%time.Second
		}
		log := fmt.Sprintf("[%s] | %s %3d %s | %s : %s  - Latency: %s - ErrorMsg: %s \n",
			params.TimeStamp.Format("Mon Jan 2 15:04:05 MST 2006"),
			statusColor, params.StatusCode, resetColor,
			params.Method,
			params.Path,
			params.Latency,
			params.ErrorMessage,
		)
		return log
	})
}*/
