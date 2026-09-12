package usecase

import (
	"context"
	"midtest/internal/auth"
	"midtest/internal/domain"
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
	S.CacheRepo.LoadToken(id)
	return &LoginOutput{Token: S.CacheRepo.GetAuth(string(id.Token))}, nil
}

func (S *Service) Delete() {}

func (S Service) EndSession(Token string) error {
	err := S.CacheRepo.DeleteItem(Token)
	return err
}
