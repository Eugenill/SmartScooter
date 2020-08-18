package middleware

import (
	"encoding/base64"
	"fmt"
	"github.com/Eugenill/SmartScooter/api_rest/models"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/db"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/sqlbunny/sqlbunny/runtime/qm"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
)

// RequestIDKey is the key that holds the unique request ID in a request context.
const RequestIDKey int = 0

var reqid uint64

func AddMiddlewares(engine *gin.Engine, ctx *gin.Context) {
	//engine.Use(RequestID())
	engine.Use(BasicAuth(ctx))
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
func BasicAuth(ctx *gin.Context) gin.HandlerFunc {
	if ctx.Request.URL.Path != "/create_user" {
		var accounts gin.Accounts
		ctx2 := db.GinToContextWithDB(ctx)
		users, err := models.Users(
			qm.Where("is_deleted = false"),
		).All(ctx2)
		if err != nil {
			errors.ErrJsonResponse(ctx, err, 401)
		}
		for _, user := range users {
			accounts[user.Username] = user.Secret
		}
		if len(accounts) != 0 {
			return gin.BasicAuth(accounts)
		}
		errors.ErrJsonResponse(ctx, errors.New("no users inserted"), 401)
	}
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}

/*
func Auth() gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		r := ctx.Request
		if r.URL.Path != "/create_user" {
			userName := r.Header.Get("username")
			secret := r.Header.Get("secret")
			if userName != "" && secret != "" {
				ctx2 := db.GinToContextWithDB(ctx)
				user, err := models.Users(
					qm.Where("username = ?", userName),
				).One(ctx2)
				if err == nil {
					if hash.CheckPasswordHash(secret, user.SecretHash) {
						ctx.Next()
					} else {
						errors.ErrJsonResponse(ctx, errors.New("Incorrect Password"), 400)
					}
				} else {
					errors.ErrJsonResponse(ctx, err, 400)
				}
			} else {
				errors.ErrJsonResponse(ctx, errors.New("Username and/or secret not inserted in header"), 400)
			}
		}
		ctx.Next()
	}
	return fn
}*/

func Logger() gin.HandlerFunc {
	fn := gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		log := fmt.Sprintf("[%s] %s : %s : %d - Latency: %s - ErrorMsg: %s \n",
			params.TimeStamp.Format("Mon Jan 2 15:04:05 MST 2006"),
			params.Method,
			params.Path,
			params.StatusCode,
			params.Latency,
			params.ErrorMessage,
		)
		return log
	})
	return fn
}
