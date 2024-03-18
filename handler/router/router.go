package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()

	healthzHandler := handler.NewHealthzHandler()

	todoService := service.NewTODOService(todoDB)
	todoHandler := handler.NewTODOHandler(todoService)

	doPanicHandler := handler.NewDoPanicHandler()
	recoverDoPanicHandler := middleware.Recovery(doPanicHandler)

	mux.Handle("/healthz", healthzHandler)
	mux.Handle("/todos", todoHandler)
	// mux.Handle("/do-panic", doPanicHandler)
	mux.Handle("/do-panic", recoverDoPanicHandler)

	return mux
}
