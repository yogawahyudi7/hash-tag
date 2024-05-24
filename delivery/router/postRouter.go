package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yogawahyudi7/social-media/delivery/controller"
	middlewares "github.com/yogawahyudi7/social-media/delivery/middleware"
)

type PostRouter struct {
	Controller *controller.PostController
	Mux        *mux.Router
	Middleware *middlewares.Middleware
}

func NewPostRouter(middleware *middlewares.Middleware, controller *controller.PostController, mux *mux.Router) *PostRouter {
	return &PostRouter{
		Controller: controller,
		Mux:        mux,
		Middleware: middleware,
	}
}

func (p *PostRouter) Register() {
	p.Mux.HandleFunc("/api/posts", p.Controller.GetAllPosts).Methods(http.MethodGet)
	p.Mux.HandleFunc("/api/posts/{post_id}", p.Controller.GetPostById).Methods(http.MethodGet)
	p.Mux.HandleFunc("/api/posts/{post_id}", p.Controller.DeletePost).Methods(http.MethodDelete)
	p.Mux.HandleFunc("/api/posts", p.Controller.CreatePost).Methods(http.MethodPost)
	p.Mux.HandleFunc("/api/posts/{post_id}", p.Controller.UpdatePost).Methods(http.MethodPut)
	p.Mux.HandleFunc("/api/posts/", p.Controller.FindByTag).Methods(http.MethodGet)
	p.Mux.HandleFunc("/api/posts/{post_id}/publish", p.Middleware.VerifyTokenMiddleware(p.Controller.PublishPost)).Methods(http.MethodPut)
}
