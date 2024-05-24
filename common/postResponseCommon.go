package common

type HttpResponse struct {
	StatusResponse StatusResponse
	PostResponse   ResponsePost
}

type StatusResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (r *StatusResponse) Success(message string, data interface{}) *StatusResponse {
	r.Code = 200
	r.Message = message
	r.Data = data
	return r
}

func (r *StatusResponse) NotFound(err error) *StatusResponse {
	r.Code = 404
	r.Message = err.Error()
	r.Data = nil
	return r
}

func (r *StatusResponse) BadRequest(message string) *StatusResponse {
	r.Code = 400
	r.Message = message
	r.Data = nil
	return r
}

func (r *StatusResponse) InternalServerError(err error) *StatusResponse {
	r.Code = 500
	r.Message = err.Error()
	r.Data = nil
	return r
}

type ResponsePost struct {
	Id      int      `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
	Status  string   `json:"status"`
}
