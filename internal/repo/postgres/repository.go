package postgres_repo

import (
	"midtest/internal/auth"
	"midtest/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	Pool *pgxpool.Pool
}

func (P *PostgresRepository) Register(Login string, Password string) bool {
	if !(auth.LoginValidate(Login) && auth.PasswordValidate(Password)) {
		return false
	}

	return true
}

func (P *PostgresRepository) LogIn(Login string, Password string) auth.AuthInfo {
	var Res auth.AuthInfo

	return Res
}
func (P *PostgresRepository) FileInfo(ID auth.ID) domain.FileInfo {
	var Res domain.FileInfo
	return Res
}
