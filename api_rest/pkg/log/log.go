package log

import (
	"context"
	"github.com/Eugenill/SmartScooter/api_rest/pkg/rest"
	"log"
)

type key int

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

func Request(ctx context.Context) {
	log.Println(ctx, rest.RequestIDKey, "")
}
