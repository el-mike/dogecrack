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

func NewServer(manager pitbull.PitbullManager, client pitbull.PitbullClient) *Server {
	router := mux.NewRouter()

	controller := NewController(manager, client)

	router.HandleFunc("/health", controller.GetHealth).Methods("GET")

	http.Handle("/", router)

	return &Server{
		router: router,
	}
}

func (s *Server) Run() {
	log.Fatal(http.ListenAndServe(":8080", s.router))
}

func NoopHandler() {}
