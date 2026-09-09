package domain

var (
	ResponseOK           = &DomainErrorField{Code: 200, Text: "OK"}
	ErrNotFound          = &DomainErrorField{Code: 405, Text: "Bad request method"}
	ErrBadParameter      = &DomainErrorField{Code: 400, Text: "Bad parameter"}
	ErrUnauthorized      = &DomainErrorField{Code: 401, Text: "Unauthorized"}
	ErrInternalServer    = &DomainErrorField{Code: 500, Text: "Internal server error"}
	ErrNoAuthForAction   = &DomainErrorField{Code: 403, Text: "Not authorized to perform this action"}
	ErrMethodNotRealized = &DomainErrorField{Code: 501, Text: "Method not realized"}
)
