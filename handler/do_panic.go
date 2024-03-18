package handler

import (
	"fmt"
	"net/http"
)

type DoPanicHandler struct {
	Count int
}

func NewDoPanicHandler() *DoPanicHandler {
	return &DoPanicHandler{}
}

func (h *DoPanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Count++
	fmt.Println("start panic:", h.Count)
	panic(fmt.Sprintf("パニーザ「私の戦闘力は%dです」", h.Count))
}
