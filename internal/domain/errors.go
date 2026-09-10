package domain

import "errors"

var ResponseOK = &DomainErrorField{Code: 200, Text: "OK"}

// http resp errors
var ErrBadRequestMethod = errors.New("bad request method")
var ErrBadParameter = errors.New("bad parameter")
var ErrUnauthorized = errors.New("unathorized")
var ErrInternal = errors.New("Internal server error")
var ErrMethodNotReady = errors.New("Method not realized")
var ErrActionNotAuthorized = errors.New("Not authorized to perform this action")
