package postgres_repo_test

import (
	"context"
	"log/slog"
	"midtest/internal/auth"
	"midtest/internal/infrastructure/postgres"
	postgres_repo "midtest/internal/repo/postgres"
	"testing"
)

var DSN = "postgresql://postgres:password@127.0.0.1:5432/test?sslmode=disable"

func TestRegister(t *testing.T) {
	ctx := context.Background()
	pool, err := postgres.Open(ctx, DSN)
	if err != nil {
		slog.Error(err.Error())
		t.Errorf("found error")
		return
	}
	Repo := postgres_repo.NewRepo(pool)
	Login, err := Repo.Register(ctx, auth.LoginInfo{Login: "DefaultLogin", Password: "121Rd@rdsdsdsd"})
	if err != nil {
		slog.Error(err.Error())
		t.Errorf("found error")
		return
	}
	if Login.Login != "DefaultLogin" {
		t.Error("login do not match")
		return
	}
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	pool, err := postgres.Open(ctx, DSN)
	if err != nil {
		slog.Error(err.Error())
		t.Errorf("found error")
		return
	}
	Repo := postgres_repo.NewRepo(pool)
	Login, err := Repo.Register(ctx, auth.LoginInfo{Login: "DefaultLogin", Password: "121Rd@rdsdsdsd"})
	if err != nil {
		t.Errorf("found error")
		return
	}
	if Login.Login != "DefaultLogin" {
		t.Error("login do not match")
		return
	}
	Res, err := Repo.LogIn(ctx, auth.LoginInfo{Login: "DefaultLogin", Password: "121Rd@rdsdsdsd"})
	if err != nil {
		t.Errorf("found error  during auth")
		return
	}
	t.Log(Res.Token)
}

func TestDelete(t *testing.T) {
	var id string
	query := `INSERT INTO Files (user_id) VALUES ($1) RETURNING id`
	ctx := context.Background()
	pool, err := postgres.Open(ctx, DSN)
	if err != nil {
		slog.Error(err.Error())
		t.Errorf("found error")
		return
	}
	Repo := postgres_repo.NewRepo(pool)
	_, err = Repo.Register(ctx, auth.LoginInfo{Login: "DefaultLogin", Password: "121Rd@rdsdsdsd"})
	if err != nil {
		t.Errorf("found error")
		return
	}
	auth, err := Repo.LogIn(ctx, auth.LoginInfo{Login: "DefaultLogin", Password: "121Rd@rdsdsdsd"})
	row := pool.QueryRow(ctx, query, auth.ID)
	err = row.Scan(&id)
	if err != nil {
		slog.Error(err.Error())
		t.Errorf("found error")
		return
	}
	err = Repo.DeleteFile(ctx, string(auth.ID))
	if err != nil {
		slog.Error(err.Error())
		t.Errorf("found error")
		return
	}
}
