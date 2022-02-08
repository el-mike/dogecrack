package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
}

func NewServer(controller *Controller) *Server {
	router := mux.NewRouter()

	router.HandleFunc("/health", controller.GetHealth).Methods("GET")
	router.HandleFunc("/getActiveInstances", controller.GetActiveInstances).Methods("GET")
	router.HandleFunc("/getInstance", controller.GetInstance).Methods("GET")

	router.HandleFunc("/getJobs", controller.GetJobs).Methods("GET")

	router.HandleFunc("/runCommand", controller.RunCommand).Methods("POST")
	router.HandleFunc("/crack", controller.Crack).Methods("POST")

	http.Handle("/", router)

	return &Server{
		router: router,
	}
}

func (s *Server) Run() {
	log.Fatal(http.ListenAndServe(":8080", s.router))
}
