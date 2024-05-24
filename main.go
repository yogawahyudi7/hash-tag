package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yogawahyudi7/social-media/common"
	"github.com/yogawahyudi7/social-media/config"
	"github.com/yogawahyudi7/social-media/delivery/controller"
	"github.com/yogawahyudi7/social-media/delivery/router"
	"github.com/yogawahyudi7/social-media/repository"
	"github.com/yogawahyudi7/social-media/util"
)

func main() {

	// setup configuration
	setup := &config.Server{}
	setup.Load()

	// connect to database
	db := util.NewDatabase(setup)

	// common
	response := &common.HttpResponse{}
	request := &common.HttpRequest{}

	// repository
	postRepository := repository.NewPostRepository(db)

	// controller
	postController := controller.NewPostController(postRepository, request, response)

	// 	initial http serve
	serve := mux.NewRouter()

	// router post registered
	postRouter := router.NewPostRouter(postController, serve)
	postRouter.Register()

	log.Printf("Server started at %v !\n", setup.AppPort)
	log.Fatal(http.ListenAndServe(setup.AppPort, serve))
}
