package router

import (
	"server/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/userinfo", controller.CreateUserInfo).Methods("POST")
	router.HandleFunc("/api/userinfo/{name}", controller.GetUserInfo).Methods("GET")
	router.HandleFunc("/api/jobinfo", controller.GetJobInfo).Methods("GET")
	return router
}
