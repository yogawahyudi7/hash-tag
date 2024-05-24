package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yogawahyudi7/social-media/common"
	"github.com/yogawahyudi7/social-media/config"
	"github.com/yogawahyudi7/social-media/delivery/controller"
	middlewares "github.com/yogawahyudi7/social-media/delivery/middleware"
	"github.com/yogawahyudi7/social-media/delivery/router"
	"github.com/yogawahyudi7/social-media/repository"
	"github.com/yogawahyudi7/social-media/util"
)

func main() {

	setup := &config.Server{}
	setup.Load()

	db := util.NewDatabase(setup)

	// common
	response := &common.HttpResponse{}
	request := &common.HttpRequest{}

	// repository
	postRepository := repository.NewPostRepository(db)
	userRepository := repository.NewUserRepository(db)

	// controller
	postController := controller.NewPostController(postRepository, request, response)
	userController := controller.NewUserController(userRepository, setup, request, response)

	// 	initial http serve
	serve := mux.NewRouter()

	// middleware
	middleware := middlewares.NewMiddleware(setup, response)

	// router post registered
	postRouter := router.NewPostRouter(middleware, postController, serve)
	postRouter.Register()

	// router user registered
	userRouter := router.NewUserRouter(userController, serve)
	userRouter.Register()

	log.Printf("Server started at %v !\n", setup.AppPort)
	log.Fatal(http.ListenAndServe(setup.AppPort, serve))
}
