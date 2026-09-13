package usecase

import (
	"context"
	"midtest/internal/auth"
	"midtest/internal/domain"
	"time"
)

type Usecase interface {
	Register(context.Context, RegisterInput) (*RegisterOutput, error)
	LogIn(ctx context.Context, AuthInfo LoginInput) (*LoginOutput, error)
	EndSession(Token string) error
}

type PostgresRepo interface {
	Register(ctx context.Context, AuthInfo auth.LoginInfo) (*auth.LoginInfo, error)
	LogIn(ctx context.Context, AuthInfo auth.LoginInfo) (*auth.AuthInfo, error)
	GetFile(ctx context.Context, ID string) (domain.FileInfo, error)
	DeleteFile(ctx context.Context, ID string) error
}

type CacheRepo interface {
	LoadFile(Name string, Data interface{}, Duration time.Duration)
	LoadAuth(Info auth.LoginInfo)
	LoadToken(AuthInfo auth.AuthInfo)
	IsActive(Token string) bool
	GetAuth(Token string) string
	CheckPWD(Info auth.LoginInfo) bool
	DeleteItem(key string) error
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
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginOutput struct {
	Token string `json:"token"`
}

type DeleteInput struct {
	ID string `json:"id"`
}

type EndSessionInput struct {
	Token string `json:"token"`
}

type EndSessionOutput struct {
	Success bool `json:"success"`
}
