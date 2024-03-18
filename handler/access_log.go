package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/model"
)

type AccessLogHandler struct{}

func NewAccessLogHandler() *AccessLogHandler {
	return &AccessLogHandler{}
}

func (h *AccessLogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	//任意の処理
	var name string
	fmt.Print("What's your name?: ")
	fmt.Scan(&name)
	fmt.Fprintf(w, "Hello %s !\n", name)
	ctx := middleware.SetUserAgent(r)

	end := time.Now()
	dif := end.Sub(start).Milliseconds()

	out := &model.AccessLog{
		Timestamp: start,
		Latency:   dif,
		Path:      r.URL.EscapedPath(),
		OS:        middleware.GetUserAgent(ctx, "OS").(string),
	}

	err := json.NewEncoder(os.Stdout).Encode(out)
	if err != nil {
		log.Println(err)
		return
	}
}
