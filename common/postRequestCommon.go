package common

type HttpRequest struct {
	CreatePostRequest CreatePostRequest
	UpdatePostRequest UpdatePostRequest
}

type CreatePostRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type UpdatePostRequest struct {
	CreatePostRequest
}
