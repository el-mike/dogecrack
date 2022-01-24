package server

import (
	"log"
	"net/http"

	"github.com/el-mike/dogecrack/shepherd/pitbull"
	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
}

func NewServer(manager *pitbull.PitbullManager) *Server {
	router := mux.NewRouter()

	controller := NewController(manager)

	router.HandleFunc("/health", controller.GetHealth).Methods("GET")
	router.HandleFunc("/getActiveInstances", controller.GetActiveInstances).Methods("GET")
	router.HandleFunc("/getInstance", controller.GetInstance).Methods("GET")
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
