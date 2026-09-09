package auth

import (
	"unicode"

	"github.com/google/uuid"
)

const PWD_MIN_LENGTH = 8
const LOGIN_MIN_LENGTH = 8

type AuthToken string
type ID string

type AuthInfo struct {
	Token AuthToken
	ID    ID
}

type LoginInfo struct {
	Login    Login    `json:"login"`
	Password Password `json:"password"`
}

func CreateToken() AuthToken {
	Token := uuid.New().String()
	return AuthToken(Token)
}

type Login string

func LoginValidate(L string) bool {
	if len(L) < LOGIN_MIN_LENGTH {
		return false
	}
	return true
}

type Password string

func PasswordValidate(P string) bool {
	if len(P) < PWD_MIN_LENGTH {
		return false
	}

	var hasLower, hasUpper, hasDigit, hasSymbol bool

	for _, r := range P {
		switch {
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r) || unicode.IsMark(r) || (!unicode.IsLetter(r) && !unicode.IsNumber(r)):
			hasSymbol = true
		}
	}
	differentCases := hasLower && hasUpper
	return differentCases && hasDigit && hasSymbol
}
