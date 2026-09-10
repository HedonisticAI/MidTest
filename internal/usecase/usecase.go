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
	S.CacheRepo.LoadAuth(*Res)
	return &RegisterOutput{Login: R.Login}, nil
}

func (S *Service) LogIn(ctx context.Context, AuthInfo auth.LoginInfo) (*LoginOutput, error) {
	info := auth.LoginInfo{Login: AuthInfo.Login, Password: AuthInfo.Password}
	id, err := S.PostgresRepo.Register(ctx, info)
	if !S.CacheRepo.CheckPWD(info) || err != nil {
		return nil, domain.ErrBadParameter
	}
	S.CacheRepo.LoadAuth(*id)
	return nil, nil
}
