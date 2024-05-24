package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/yogawahyudi7/social-media/common"
	"github.com/yogawahyudi7/social-media/constant"
	"github.com/yogawahyudi7/social-media/model"
	"github.com/yogawahyudi7/social-media/repository"
)

type PostController struct {
	PostRepository *repository.PostRepository
	Response       *common.HttpResponse
	Request        *common.HttpRequest
}

func NewPostController(repository *repository.PostRepository, request *common.HttpRequest, response *common.HttpResponse) *PostController {
	return &PostController{
		PostRepository: repository,
		Response:       response,
		Request:        request,
	}
}

func (pc *PostController) GetAllPosts(w http.ResponseWriter, r *http.Request) {

	w.Header().Set(constant.ContentType, constant.ApplicationJSON)

	posts, err := pc.PostRepository.FindAllWithTag()
	if err != nil {
		json.NewEncoder(w).Encode(pc.Response.StatusResponse.NotFound(err))
		return
	}

	if len(posts) < 1 {
		json.NewEncoder(w).Encode(pc.Response.StatusResponse.NotFound(errors.New(constant.ErrorNotFound)))
		return
	}

	postsResponse := []common.ResponsePost{}
	for _, post := range posts {
		tags := []string{}
		for _, tag := range post.Tags {
			tags = append(tags, tag.Label)
		}

		postsResponse = append(postsResponse, common.ResponsePost{
			Id:      post.ID,
			Title:   post.Title,
			Content: post.Content,
			Tags:    tags,
			Status:  post.Status,
		})
	}

	json.NewEncoder(w).Encode(pc.Response.StatusResponse.Success(constant.DataFound, postsResponse))

}

func (pc *PostController) GetPostById(w http.ResponseWriter, r *http.Request) {

	w.Header().Set(constant.ContentType, constant.ApplicationJSON)

	idStr := mux.Vars(r)["post_id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		json.NewEncoder(w).Encode(pc.Response.StatusResponse.BadRequest(constant.ErrorBadRequest))
		return
	}

	post, err := pc.PostRepository.FindByID(id)
	if err != nil {
		json.NewEncoder(w).Encode(pc.Response.StatusResponse.NotFound(err))
		return
	}

	if post.ID == 0 {
		json.NewEncoder(w).Encode(pc.Response.StatusResponse.NotFound(errors.New(constant.ErrorNotFound)))
		return
	}

	tags := []string{}
	for _, tag := range post.Tags {
		tags = append(tags, tag.Label)
	}

	postsResponse := common.ResponsePost{
		Id:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		Tags:    tags,
		Status:  post.Status,
	}

	json.NewEncoder(w).Encode(pc.Response.StatusResponse.Success(constant.DataFound, postsResponse))
}

func (pc *PostController) DeletePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constant.ContentType, constant.ApplicationJSON)

	idStr := mux.Vars(r)["post_id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		json.NewEncoder(w).Encode(pc.Response.StatusResponse.BadRequest(constant.ErrorBadRequest))
		return
	}

	rows, err := pc.PostRepository.HardDelete(id)
	if err != nil {
		json.NewEncoder(w).Encode(pc.Response.StatusResponse.InternalServerError(err))
		return
	}

	if rows < 1 {
		json.NewEncoder(w).Encode(pc.Response.StatusResponse.NotFound(errors.New(constant.ErrorNotFound)))
		return
	}

	json.NewEncoder(w).Encode(pc.Response.StatusResponse.Success(constant.DataDeleted, nil))
}

func (pc *PostController) CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constant.ContentType, constant.ApplicationJSON)

	postRequest := pc.Request.CreatePostRequest
	err := json.NewDecoder(r.Body).Decode(&postRequest)
	if err != nil {
		json.NewEncoder(w).Encode(pc.Response.StatusResponse.BadRequest(constant.ErrorBadRequest))
		return
	}

	postRequestStatus := constant.PostStatusDraft
	post := model.Post{
		Title:   postRequest.Title,
		Content: postRequest.Content,
		Status:  postRequestStatus,
	}

	tags := []model.Tag{}
	for _, tag := range postRequest.Tags {
		tags = append(tags, model.Tag{Label: strings.ToLower(tag)})
	}

	postId, err := pc.PostRepository.Create(&post, &tags)
	if err != nil {
		json.NewEncoder(w).Encode(pc.Response.StatusResponse.InternalServerError(err))
		return
	}

	if postId < 1 {
		json.NewEncoder(w).Encode(pc.Response.StatusResponse.NotFound(errors.New(constant.ErrorNotFound)))
		return
	}

	dataResponse := map[string]int{"id": postId}

	json.NewEncoder(w).Encode(pc.Response.StatusResponse.Success(constant.DataSaved, dataResponse))
}

func (pc *PostController) UpdatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constant.ContentType, constant.ApplicationJSON)

	postRequest := pc.Request.UpdatePostRequest
	err := json.NewDecoder(r.Body).Decode(&postRequest)
	if err != nil {
		json.NewEncoder(w).Encode(pc.Response.StatusResponse.BadRequest(constant.ErrorBadRequest))
		return
	}

	idStr := mux.Vars(r)["post_id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		json.NewEncoder(w).Encode(pc.Response.StatusResponse.BadRequest(constant.ErrorBadRequest))
		return
	}

	postRequestStatus := constant.PostStatusDraft
	post := model.Post{
		ID:      id,
		Title:   postRequest.Title,
		Content: postRequest.Content,
		Status:  postRequestStatus,
	}

	tags := []model.Tag{}
	for _, tag := range postRequest.Tags {
		tags = append(tags, model.Tag{Label: strings.ToLower(tag)})
	}

	id, err = pc.PostRepository.Update(&post, &tags)
	if err != nil {
		json.NewEncoder(w).Encode(pc.Response.StatusResponse.InternalServerError(err))
		return
	}

	if id < 1 {
		json.NewEncoder(w).Encode(pc.Response.StatusResponse.NotFound(errors.New(constant.ErrorNotFound)))
		return
	}

	dataResponse := map[string]int{"id": id}
	json.NewEncoder(w).Encode(pc.Response.StatusResponse.Success(constant.DataUpdated, dataResponse))
}

func (pc *PostController) FindByTag(w http.ResponseWriter, r *http.Request) {

	w.Header().Set(constant.ContentType, constant.ApplicationJSON)

	queryParams := r.URL.Query()

	label := queryParams.Get("tag")

	if label == "" {
		json.NewEncoder(w).Encode(pc.Response.StatusResponse.BadRequest(constant.ErrorBadRequest))
		return
	}

	ids, err := pc.PostRepository.FindIdByTag(label)
	if err != nil {
		json.NewEncoder(w).Encode(pc.Response.StatusResponse.InternalServerError(err))
		return
	}

	if len(ids) < 1 {
		json.NewEncoder(w).Encode(pc.Response.StatusResponse.NotFound(errors.New(constant.ErrorNotFound)))
		return
	}
	posts, err := pc.PostRepository.FindByIds(ids)
	if err != nil {
		json.NewEncoder(w).Encode(pc.Response.StatusResponse.InternalServerError(err))
		return
	}

	if len(posts) < 1 {
		json.NewEncoder(w).Encode(pc.Response.StatusResponse.NotFound(errors.New(constant.ErrorNotFound)))
		return
	}

	postsResponse := []common.ResponsePost{}
	for _, post := range posts {
		tags := []string{}
		for _, tag := range post.Tags {
			tags = append(tags, tag.Label)
		}

		postsResponse = append(postsResponse, common.ResponsePost{
			Id:      post.ID,
			Title:   post.Title,
			Content: post.Content,
			Tags:    tags,
			Status:  post.Status,
		})
	}

	json.NewEncoder(w).Encode(pc.Response.StatusResponse.Success(constant.DataFound, postsResponse))
}
