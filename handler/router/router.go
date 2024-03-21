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
	authHealthzHandler := middleware.BasicAuthMiddleware(healthzHandler)

	todoService := service.NewTODOService(todoDB)
	todoHandler := handler.NewTODOHandler(todoService)

	doPanicHandler := handler.NewDoPanicHandler()
	recoverDoPanicHandler := middleware.Recovery(doPanicHandler)

	interactiveHandler := handler.NewInteractiveHandler()
	authInteractiveHandler := middleware.BasicAuthMiddleware(interactiveHandler)

	getAccessLogHandler := middleware.AccessLog(interactiveHandler)
	accessLogHandler := middleware.SetOSMiddleware(getAccessLogHandler)

	mux.Handle("/healthz", authHealthzHandler)
	mux.Handle("/todos", todoHandler)
	mux.Handle("/do-panic", recoverDoPanicHandler)
	mux.Handle("/access_log", accessLogHandler)
	mux.Handle("/interactive", authInteractiveHandler)

	return mux
}
