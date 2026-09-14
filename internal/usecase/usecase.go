package usecase

import (
	"context"
	"midtest/internal/auth"
	"midtest/internal/domain"
	"os"
	"slices"
	"time"
)

type Service struct {
	ADMToken     string
	PostgresRepo PostgresRepo
	CacheRepo    CacheRepo
}

func NewUsecase(ADMToken string, PostgresRepo PostgresRepo, CacheRepo CacheRepo) Usecase {
	return &Service{
		ADMToken:     ADMToken,
		PostgresRepo: PostgresRepo,
		CacheRepo:    CacheRepo,
	}
}

func (S *Service) Register(ctx context.Context, R RegisterInput) (*RegisterOutput, error) {
	if R.AdmToken != S.ADMToken {
		return nil, domain.ErrActionNotAuthorized
	}
	Auth := auth.LoginInfo{Login: R.Login, Password: R.Password}
	Res, err := S.PostgresRepo.Register(ctx, Auth)
	if err != nil {
		return nil, err
	}
	S.CacheRepo.LoadAuth(Auth)
	return &RegisterOutput{Login: Res.Login}, nil
}

func (S *Service) LogIn(ctx context.Context, AuthInfo LoginInput) (*LoginOutput, error) {
	info := auth.LoginInfo{Login: AuthInfo.Login, Password: AuthInfo.Password}
	id, err := S.PostgresRepo.LogIn(ctx, info)
	S.CacheRepo.LoadAuth(info)
	if !S.CacheRepo.CheckPWD(info) || err != nil {
		return nil, domain.ErrBadParameter
	}
	token := auth.AuthInfo{ID: id.ID, Token: id.Token}
	S.CacheRepo.LoadToken(token)
	return &LoginOutput{Token: S.CacheRepo.GetAuth(string(id.Token))}, nil
}

func (S *Service) EndSession(Token string) error {
	err := S.CacheRepo.DeleteItem(Token)
	return err
}

func (S *Service) ListFiles(ctx context.Context, List ListInput) (interface{}, error) {
	Data, err := S.PostgresRepo.List(ctx, List.Filters)
	if err != nil {
		return nil, err
	}
	return Data, nil
}

func (S *Service) GetFile(ctx context.Context, ID string, Token string) (interface{}, error) {
	if !S.CacheRepo.IsActive(Token) {
		return nil, domain.ErrBadParameter
	}
	Res, err := S.PostgresRepo.GetFileInfo(ctx, ID)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(Res.Users, ID) {
		return nil, domain.ErrActionNotAuthorized
	}
	file, exists := S.CacheRepo.GetFile(ID)
	if exists {
		return file, nil
	}
	memfile, err := os.ReadFile(Res.Path + Res.Name)
	if err != nil {
		return nil, err
	}
	S.CacheRepo.LoadFile(Res.Name, memfile, 10*time.Minute)
	return memfile, nil
}

func (S *Service) DeleteFile(ctx context.Context, ID string, Token string) error {
	if !S.CacheRepo.IsActive(Token) {
		return domain.ErrBadParameter
	}
	Res, err := S.PostgresRepo.GetFileInfo(ctx, ID)
	if err != nil {
		return err
	}
	if !slices.Contains(Res.Users, ID) {
		return domain.ErrActionNotAuthorized
	}
	err = S.PostgresRepo.DeleteFile(ctx, ID)
	if err != nil {
		return err
	}
	S.CacheRepo.DeleteItem(ID)
	return nil
}
