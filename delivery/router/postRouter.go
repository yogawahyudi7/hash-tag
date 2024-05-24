package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yogawahyudi7/social-media/delivery/controller"
)

type PostRouter struct {
	Controller *controller.PostController
	Mux        *mux.Router
}

func NewPostRouter(controller *controller.PostController, mux *mux.Router) *PostRouter {
	return &PostRouter{
		Controller: controller,
		Mux:        mux,
	}
}

func (p *PostRouter) Register() {
	p.Mux.HandleFunc("/api/posts", p.Controller.GetAllPosts).Methods(http.MethodGet)
	p.Mux.HandleFunc("/api/posts/{post_id}", p.Controller.GetPostById).Methods(http.MethodGet)
	p.Mux.HandleFunc("/api/posts/{post_id}", p.Controller.DeletePost).Methods(http.MethodDelete)
	p.Mux.HandleFunc("/api/posts", p.Controller.CreatePost).Methods(http.MethodPost)
	p.Mux.HandleFunc("/api/posts/{post_id}", p.Controller.UpdatePost).Methods(http.MethodPut)
	p.Mux.HandleFunc("/api/posts/", p.Controller.FindByTag).Methods(http.MethodGet)
}
