package usecase_test

import (
	"context"
	"log/slog"
	"midtest/internal/infrastructure/cache"
	"midtest/internal/infrastructure/postgres"
	cache_repo "midtest/internal/repo/cache"
	postgres_repo "midtest/internal/repo/postgres"
	"midtest/internal/usecase"
	"testing"
	"time"
)

var DSN = "postgresql://postgres:password@127.0.0.1:5432/test?sslmode=disable"

func CreateUsecase(Cache *cache_repo.RepoCache) (usecase.Usecase, error) {
	ctx := context.Background()
	pool, err := postgres.Open(ctx, DSN)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	Repo := postgres_repo.NewRepo(pool)
	return usecase.NewUsecase("TESTTOKEN", Repo, Cache), nil
}

func TestRegister(t *testing.T) {
	cacheInstance := cache.NewCache(5*time.Minute, 10*time.Minute)
	RepoCache := &cache_repo.RepoCache{
		Cache: cacheInstance,
	}
	Usecase, err := CreateUsecase(RepoCache)
	if err != nil {
		t.Errorf("found error during usecase creation")
		return
	}
	Data := usecase.RegisterInput{AdmToken: "TESTTOKEN", Login: "DefaultLogin", Password: "121Rd@rdsdsdsd"}
	Res, err := Usecase.Register(context.Background(), Data)
	if err != nil {
		t.Errorf("found error during register")
		return
	}
	if Res.Login != "DefaultLogin" {
		t.Error("login do not match")
		return
	}
}

func TestLogin(t *testing.T) {
	cacheInstance := cache.NewCache(5*time.Minute, 10*time.Minute)
	RepoCache := &cache_repo.RepoCache{
		Cache: cacheInstance,
	}
	Usecase, err := CreateUsecase(RepoCache)
	if err != nil {
		t.Errorf("found error during usecase creation")
		return
	}
	Data := usecase.LoginInput{Login: "DefaultLogin", Password: "121Rd@rdsdsdsd"}
	Res, err := Usecase.LogIn(context.Background(), Data)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("found error during login")
		return
	}
	if Res.Token == "" {
		t.Error("token is empty")
		return
	}
}

func TestEndSession(t *testing.T) {
	cacheInstance := cache.NewCache(5*time.Minute, 10*time.Minute)
	RepoCache := &cache_repo.RepoCache{
		Cache: cacheInstance,
	}
	Usecase, err := CreateUsecase(RepoCache)
	if err != nil {
		t.Errorf("found error during usecase creation")
		return
	}
	Data := usecase.LoginInput{Login: "DefaultLogin", Password: "121Rd@rdsdsdsd"}
	Res, err := Usecase.LogIn(context.Background(), Data)
	t.Log(Res.Token)
	if err != nil {
		t.Log(err.Error())
		t.Errorf("found error during login")
		return
	}
	t.Log(RepoCache.IsActive(Res.Token))
	if Res.Token == "" {
		t.Error("token is empty")
		return
	}
	t.Log("Ending session...")
	EndSessionRes := Usecase.EndSession(Res.Token)
	if EndSessionRes != nil {
		t.Error(EndSessionRes)
		return
	}
}
