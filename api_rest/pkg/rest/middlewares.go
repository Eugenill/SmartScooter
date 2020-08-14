package rest

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/Eugenill/SmartScooter/api_rest/models"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/db"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/hash"
	"github.com/go-chi/chi"
	"github.com/sqlbunny/sqlbunny/runtime/bunny"
	"github.com/sqlbunny/sqlbunny/runtime/qm"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
)

// Key to use when setting the request ID.
type ctxKeyRequestID int

// RequestIDKey is the key that holds the unique request ID in a request context.
const RequestIDKey ctxKeyRequestID = 0

var reqid uint64

func AddMiddlewares(router *chi.Mux) {
	router.Use(RequestID)
	router.Use(CtxWithDB)
	router.Use(Auth)
}

// RequestID is a middleware that injects a request ID into the context of each
// request. A request ID is a string of the form "host.example.com/random-0001",
// where "random" is a base62 random string that uniquely identifies this go
// process, and where the last number is an atomically incremented request
// counter.
func RequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := r.Header.Get("X-Request-Id")
		if requestID == "" {
			myid := atomic.AddUint64(&reqid, 1)
			requestID = fmt.Sprintf("%s-%06d", prefixGen(), myid)
		}
		ctx = context.WithValue(ctx, RequestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
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

func CtxWithDB(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := bunny.ContextWithDB(r.Context(), db.DB)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

type userLogin struct {
	Login  string `json:"username" `
	Secret string `json:"secret" `
}

func Auth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		usr := userLogin{}
		if err := UnmarshalJSONRequest(&usr, r); err == nil {
			user, err := models.Users(
				qm.Where("login = ?", usr.Login),
			).One(ctx)
			if err == nil {
				if hash.CheckPasswordHash(usr.Secret, user.SecretHash) {
					_, err = w.Write([]byte("Correct User and password"))
					next.ServeHTTP(w, r.WithContext(ctx))
				} else {
					_, err = w.Write([]byte("Incorrect password"))
					next.ServeHTTP(w, r.WithContext(ctx))
					log.Fatal(err)
				}
			} else {
				_, err = w.Write([]byte("This user doesn't exist"))
				next.ServeHTTP(w, r.WithContext(ctx))
				log.Fatal(err)
			}
		} else {
			_, err = w.Write([]byte("Incorrect syntax"))
			next.ServeHTTP(w, r.WithContext(ctx))
			log.Fatal(err)
		}
	}
	return http.HandlerFunc(fn)
}
