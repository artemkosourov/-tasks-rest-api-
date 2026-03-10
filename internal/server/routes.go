package server

import (
	"net/http"
	"tasks-rest-api/internal/middleware"
	"tasks-rest-api/internal/service"
)

func SetupRoutes(mux *http.ServeMux, svc *service.TaskService) {
	mux.HandleFunc("POST /register", http.HandlerFunc(RegisterHandler(svc)))
	mux.HandleFunc("POST /login", http.HandlerFunc(LoginHandler(svc)))
	mux.HandleFunc("POST /tasks", middleware.LoggerMiddleware(middleware.AuthMiddleware(CreateTaskHandler(svc))))
	mux.HandleFunc("GET /tasks/", middleware.LoggerMiddleware(middleware.AuthMiddleware(GetTaskHandler(svc))))
	mux.HandleFunc("GET /external/users", middleware.LoggerMiddleware(middleware.AuthMiddleware(GetExternalUsersHandler(svc))))
}
