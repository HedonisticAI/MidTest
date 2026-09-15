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
	//GetFile(ctx context.Context, ID string, Token string) (interface{}, error)
	ListFiles(ctx context.Context, List ListInput) ([]domain.FileInfo, error)
	DeleteFile(ctx context.Context, ID string, Token string) error
}

type PostgresRepo interface {
	Register(ctx context.Context, AuthInfo auth.LoginInfo) (*auth.LoginInfo, error)
	LogIn(ctx context.Context, AuthInfo auth.LoginInfo) (*auth.AuthInfo, error)
	GetFileInfo(ctx context.Context, ID string) (domain.FileInfo, error)
	DeleteFile(ctx context.Context, ID string) error
	List(ctx context.Context, Filters map[string]string) ([]domain.FileInfo, error)
}

type CacheRepo interface {
	LoadFile(Name string, Data interface{}, Duration time.Duration)
	LoadAuth(Info auth.LoginInfo)
	LoadToken(AuthInfo auth.AuthInfo)
	IsActive(Token string) bool
	GetAuth(Token string) string
	CheckPWD(Info auth.LoginInfo) bool
	DeleteItem(key string) error
	GetFile(ID string) (interface{}, bool)
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

type ListInput struct {
	Filters map[string]string `json:"filters"`
}

type WriteFileInput struct {
	Name  string   `json:"name"`
	File  bool     `json:"file"`
	Token string   `json:"token"`
	Users []string `json:"grant"`
}

type WriteFileOutput struct {
	Json []byte `json:"json,omitempty"`
	Name string `json:"name"`
}
