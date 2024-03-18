package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/model"
)

// A HealthzHandler implements health check endpoint.
type HealthzHandler struct {
}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewHealthzHandler() *HealthzHandler {
	return &HealthzHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := &model.HealthzResponse{Message: "OK"}
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Println(err)
	}

	ctx := middleware.SetUserAgent(r)
	fmt.Println("Client OS:", middleware.GetUserAgent(ctx, "OS"))
}
