package log

import (
	"context"
	"log"
	"math/rand"
	"net/http"
)

type key int

const requestIDKey = key(50)

//function to put the msg to the key in the context
func Println(ctx context.Context, key key, msg string) {
	id, ok := ctx.Value(key).(int64)
	if key == 50 {
		id, ok = ctx.Value(key).(int64)
	}
	if !ok {
		log.Println("could not find any value in the context key")
		return
	}
	log.Printf("[%d] %s", id, msg)
}

func AddReqID(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := rand.Int63()
		ctx = context.WithValue(ctx, requestIDKey, id)
		f(w, r.WithContext(ctx))
	}
}
