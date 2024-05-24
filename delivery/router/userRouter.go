package router

import (
	"github.com/gorilla/mux"
	"github.com/yogawahyudi7/hash-tag/delivery/controller"
)

type UserRouter struct {
	Controller *controller.UserController
	Mux        *mux.Router
}

func NewUserRouter(controller *controller.UserController, mux *mux.Router) *UserRouter {
	return &UserRouter{
		Controller: controller,
		Mux:        mux,
	}
}

func (u *UserRouter) Register() {
	u.Mux.HandleFunc("/api/register", u.Controller.Register).Methods("POST")
	u.Mux.HandleFunc("/api/login", u.Controller.Login).Methods("POST")
}
