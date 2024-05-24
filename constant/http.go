package constant

const (
	ContentType     = "Content-Type"
	ApplicationJSON = "application/json"

	ErrorNotFound     = "data not found"
	ErrorBadRequest   = "bad request"
	ErrorUnauthorized = "token unauthorized"
	ErrorForbidden    = "forbidden access"
	ErrorTokenExpired = "token expired"
	ErrorExpired      = "expired"

	PostStatusDraft     = "draft"
	PostStatusPublished = "published"

	DataSaved         = "data saved"
	DataUpdated       = "data updated"
	DataDeleted       = "data deleted"
	DataFound         = "data found"
	DataPostPublished = "post published"

	UsernameOrPasswordCannotBeEmpty = "username or password cannot be empty"
	UsernameOrPasswordIsIncorrect   = "username or password is incorrect"
	ErrorDuplicateKey               = "duplicate key"
	UsernameAlreadyExists           = "username already exists"
)
