package usecase

import (
	"context"
	"midtest/internal/auth"
	"time"
)

type Usecase interface {
	Register(context.Context, RegisterInput) (*RegisterOutput, error)
	LogIn(ctx context.Context, AuthInfo auth.LoginInfo) (*LoginOutput, error)
}

type PostgresRepo interface {
	Register(ctx context.Context, AuthInfo auth.LoginInfo) (*auth.AuthInfo, error)
	LogIn(ctx context.Context, Login string, Password string) (*auth.AuthInfo, error)
	DeleteFile(ctx context.Context, ID string) error
}

type CacheRepo interface {
	LoadFile(Name string, Data interface{}, Duration time.Duration)
	LoadAuth(AuthInfo auth.AuthInfo)
	IsActive(Token string) bool
	GetAuth(Token string) string
	CheckPWD(Info auth.LoginInfo) bool
}

type RegisterInput struct {
	AdmToken string `json:"adm_token"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegisterOutput struct {
	Login string `json:"login"`
}

type LoginInput struct {
	Login     string `json:"login"`
	Passwordd string `json:"password"`
}

type LoginOutput struct {
	Token string `json:"token"`
}
