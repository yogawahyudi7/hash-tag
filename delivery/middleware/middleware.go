package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/yogawahyudi7/social-media/common"
	"github.com/yogawahyudi7/social-media/config"
	"github.com/yogawahyudi7/social-media/constant"
	"github.com/yogawahyudi7/social-media/helper"
)

type Middleware struct {
	Config   *config.Server
	Response *common.HttpResponse
}

func NewMiddleware(config *config.Server, response *common.HttpResponse) *Middleware {
	return &Middleware{
		Config:   config,
		Response: response,
	}
}

func (m *Middleware) VerifyTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set(constant.ContentType, constant.ApplicationJSON)

		bearerToken := r.Header.Get("Authorization")

		if !strings.Contains(bearerToken, "Bearer") {
			json.NewEncoder(w).Encode(m.Response.StatusResponse.Unauthorized(errors.New(constant.ErrorUnauthorized)))
			return
		}

		token := strings.ReplaceAll(bearerToken, "Bearer ", "")

		claims, err := helper.VerifyToken(m.Config, token)
		if err != nil {
			if strings.Contains(err.Error(), constant.ErrorExpired) {
				json.NewEncoder(w).Encode(m.Response.StatusResponse.Unauthorized(errors.New(constant.ErrorTokenExpired)))
				return
			}

			json.NewEncoder(w).Encode(m.Response.StatusResponse.Unauthorized(errors.New(constant.ErrorUnauthorized)))
			return
		}

		isAdmin := claims["isAdmin"].(bool)
		if !isAdmin {
			json.NewEncoder(w).Encode(m.Response.StatusResponse.Forbidden(errors.New(constant.ErrorForbidden)))
			return
		}

		next.ServeHTTP(w, r)
	}
}
