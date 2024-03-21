package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
)

func AccessLog(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		end := time.Now()
		dif := end.Sub(start).Milliseconds()

		out := &model.AccessLog{
			Timestamp: start,
			Latency:   dif,
			Path:      r.URL.EscapedPath(),
			OS:        GetOS(r.Context()),
		}

		err := json.NewEncoder(os.Stdout).Encode(out)
		if err != nil {
			log.Println(err)
			return
		}
	}
	return http.HandlerFunc(fn)
}
