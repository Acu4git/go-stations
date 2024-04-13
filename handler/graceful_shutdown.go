package handler

import (
	"fmt"
	"net/http"
	"time"
)

type GracefulHandler struct{}

func NewGracefulHandler() *GracefulHandler {
	return &GracefulHandler{}
}

func (h *GracefulHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	fmt.Fprintln(w, "Welcome to Graceful test!")
}
