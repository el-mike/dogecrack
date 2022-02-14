package server

import (
	"log"
	"net/http"

	"github.com/el-mike/dogecrack/shepherd/internal/auth"
	"github.com/el-mike/dogecrack/shepherd/internal/pitbull"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Server struct {
	port          string
	originAllowed string

	router *mux.Router
}

func NewServer(port, originAllowed string) *Server {
	appController := NewController()

	authController := auth.NewController()
	authMiddleware := auth.NewMiddleware()

	pitbullController := pitbull.NewController()

	baseRouter := mux.NewRouter()

	baseRouter.HandleFunc("/health", appController.GetHealth).Methods("GET")

	baseRouter.HandleFunc("/login", authController.Login).Methods(http.MethodPost, http.MethodOptions)

	protectedRouter := baseRouter.PathPrefix("/").Subrouter()

	protectedRouter.Use(authMiddleware.Middleware)

	protectedRouter.HandleFunc("/me", authController.Me).Methods("GET", "OPTIONS")
	protectedRouter.HandleFunc("/logout", authController.Logout).Methods("GET")

	protectedRouter.HandleFunc("/getActiveInstances", pitbullController.GetActiveInstances).Methods("GET")
	protectedRouter.HandleFunc("/getInstance", pitbullController.GetInstance).Methods("GET")

	protectedRouter.HandleFunc("/getJobs", pitbullController.GetJobs).Methods("GET")

	protectedRouter.HandleFunc("/runCommand", pitbullController.RunCommand).Methods("POST")
	protectedRouter.HandleFunc("/crack", pitbullController.Crack).Methods("POST")

	http.Handle("/", baseRouter)

	return &Server{
		port:          port,
		originAllowed: originAllowed,
		router:        baseRouter,
	}
}

func (s *Server) Run() {
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	credentials := handlers.AllowCredentials()
	methods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	ttl := handlers.MaxAge(3600)
	origins := handlers.AllowedOrigins([]string{s.originAllowed})

	router := handlers.CORS(headers, credentials, methods, origins, ttl)(s.router)

	log.Fatal(http.ListenAndServe(":"+s.port, router))
}
