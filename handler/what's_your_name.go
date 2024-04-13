package handler

import (
	"fmt"
	"net/http"
	"strings"
)

type InteractiveHandler struct{}

func NewInteractiveHandler() *InteractiveHandler {
	return &InteractiveHandler{}
}

func (h *InteractiveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var name string
	var num int
	fmt.Fprint(w, "What's your name?: ")
	fmt.Scan(&name)
	fmt.Fprint(w, "Tell me a number you like: ")
	fmt.Scan(&num)
	fmt.Fprint(w, strings.Repeat(fmt.Sprintf("Hello %s !\n", name), num))
}
