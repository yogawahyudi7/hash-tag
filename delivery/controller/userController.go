package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/yogawahyudi7/hash-tag/common"
	"github.com/yogawahyudi7/hash-tag/config"
	"github.com/yogawahyudi7/hash-tag/constant"
	"github.com/yogawahyudi7/hash-tag/helper"
	"github.com/yogawahyudi7/hash-tag/repository"
)

type UserController struct {
	UserRepository *repository.UserRepository
	Config         *config.Server
	Response       *common.HttpResponse
	Request        *common.HttpRequest
}

func NewUserController(repository *repository.UserRepository, config *config.Server, request *common.HttpRequest, response *common.HttpResponse) *UserController {
	return &UserController{
		UserRepository: repository,
		Config:         config,
		Response:       response,
		Request:        request,
	}
}

func (uc *UserController) Register(w http.ResponseWriter, r *http.Request) {

	w.Header().Set(constant.ContentType, constant.ApplicationJSON)

	userRegister := uc.Request.UserRegisterRequest
	err := json.NewDecoder(r.Body).Decode(&userRegister)
	if err != nil {
		json.NewEncoder(w).Encode(uc.Response.StatusResponse.BadRequest(constant.ErrorBadRequest))
		return
	}

	if userRegister.Username == "" || userRegister.Password == "" {
		json.NewEncoder(w).Encode(uc.Response.StatusResponse.BadRequest(constant.UsernameOrPasswordCannotBeEmpty))
		return
	}

	hashPassword, err := helper.HashPassword(userRegister.Password)
	if err != nil {
		json.NewEncoder(w).Encode(uc.Response.StatusResponse.InternalServerError(err))
		return
	}

	_, err = uc.UserRepository.Register(userRegister.Username, hashPassword)
	if err != nil {

		if strings.Contains(err.Error(), constant.ErrorDuplicateKey) {
			json.NewEncoder(w).Encode(uc.Response.StatusResponse.BadRequest(constant.UsernameAlreadyExists))
			return
		}

		json.NewEncoder(w).Encode(uc.Response.StatusResponse.InternalServerError(err))
		return
	}

	responseData := map[string]interface{}{
		"username": userRegister.Username,
	}
	json.NewEncoder(w).Encode(uc.Response.StatusResponse.Success(constant.DataSaved, responseData))
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set(constant.ContentType, constant.ApplicationJSON)

	userLogin := uc.Request.UserRegisterRequest
	err := json.NewDecoder(r.Body).Decode(&userLogin)
	if err != nil {
		json.NewEncoder(w).Encode(uc.Response.StatusResponse.BadRequest(constant.ErrorBadRequest))
		return
	}
	if userLogin.Username == "" || userLogin.Password == "" {
		json.NewEncoder(w).Encode(uc.Response.StatusResponse.BadRequest(constant.UsernameOrPasswordCannotBeEmpty))
		return
	}

	user, err := uc.UserRepository.FindUsername(userLogin.Username)
	if err != nil {
		json.NewEncoder(w).Encode(uc.Response.StatusResponse.InternalServerError(err))
		return
	}

	if user.ID == 0 {
		json.NewEncoder(w).Encode(uc.Response.StatusResponse.Unauthorized(errors.New(constant.UsernameOrPasswordIsIncorrect)))
		return
	}

	if !helper.ComparePassword(userLogin.Password, user.Password) {
		json.NewEncoder(w).Encode((uc.Response.StatusResponse.Unauthorized(errors.New(constant.UsernameOrPasswordIsIncorrect))))
		return
	}

	generatedToken, err := helper.GenerateToken(uc.Config, user.ID, user.IsAdmin)
	if err != nil {
		json.NewEncoder(w).Encode(uc.Response.StatusResponse.InternalServerError(err))
		return
	}

	responseData := map[string]interface{}{
		"token": generatedToken,
	}
	json.NewEncoder(w).Encode(uc.Response.StatusResponse.Success(constant.DataFound, responseData))
}
