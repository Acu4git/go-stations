package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/mileusna/useragent"
)

type ctxKey struct{}

var key = ctxKey{}

func SetOSMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ua := useragent.Parse(r.UserAgent())
		ctx := r.Context()
		r = r.WithContext(context.WithValue(ctx, key, ua.OS))
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func GetOS(ctx context.Context) string {
	v, ok := ctx.Value(key).(string)
	if !ok {
		log.Println("interface{} is nil, not string")
		return ""
	}
	return v
}
