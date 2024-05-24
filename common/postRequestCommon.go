package common

type HttpRequest struct {
	CreatePostRequest CreatePostRequest
	UpdatePostRequest UpdatePostRequest

	UserRegisterRequest UserRegisterRequest
}

type CreatePostRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type UpdatePostRequest struct {
	CreatePostRequest
}

type UserRegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
