package server

import (
	"log"
	"net/http"

	"github.com/el-mike/dogecrack/shepherd/internal/auth"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull"
	"github.com/gorilla/mux"
)

type Server struct {
	port string

	router *mux.Router
}

func NewServer(port string) *Server {
	appController := NewController()

	authController := auth.NewController()
	authMiddleware := auth.NewMiddleware()

	pitbullController := pitbull.NewController()

	baseRouter := mux.NewRouter()

	baseRouter.HandleFunc("/health", appController.GetHealth).Methods("GET")

	baseRouter.HandleFunc("/login", authController.Login).Methods("POST")

	pitbullRouter := baseRouter.PathPrefix("/").Subrouter()

	pitbullRouter.Use(authMiddleware.Middleware)

	pitbullRouter.HandleFunc("/getActiveInstances", pitbullController.GetActiveInstances).Methods("GET")
	pitbullRouter.HandleFunc("/getInstance", pitbullController.GetInstance).Methods("GET")

	pitbullRouter.HandleFunc("/getJobs", pitbullController.GetJobs).Methods("GET")

	pitbullRouter.HandleFunc("/runCommand", pitbullController.RunCommand).Methods("POST")
	pitbullRouter.HandleFunc("/crack", pitbullController.Crack).Methods("POST")

	http.Handle("/", baseRouter)

	return &Server{
		port:   port,
		router: baseRouter,
	}
}

func (s *Server) Run() {
	log.Fatal(http.ListenAndServe(":"+s.port, s.router))
}
