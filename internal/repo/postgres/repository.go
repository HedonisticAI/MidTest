package postgres_repo

import (
	"context"
	"midtest/internal/auth"
	"midtest/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	Pool *pgxpool.Pool
}

func NewRepo(Pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{Pool: Pool}
}

func (P *PostgresRepository) Register(ctx context.Context, AuthInfo auth.LoginInfo) (*auth.LoginInfo, error) {

	const query = `INSERT INTO Users (login, password) VALUES ($1, $2) RETURNING login, password`
	var Res auth.LoginInfo
	if !(auth.LoginValidate(AuthInfo.Login) && auth.PasswordValidate(AuthInfo.Password)) {
		return nil, ErrBadAuth
	}
	row := P.Pool.QueryRow(ctx, query, AuthInfo.Login, AuthInfo.Password)
	err := row.Scan(&Res.Login, &Res.Password)
	if err != nil {
		return nil, err
	}
	return &Res, nil
}

func (P *PostgresRepository) LogIn(ctx context.Context, AuthInfo auth.LoginInfo) (*auth.AuthInfo, error) {
	var Res auth.AuthInfo
	const query = `	SELECT id FROM Users WHERE login = $1 AND password = $2`
	row := P.Pool.QueryRow(ctx, query, AuthInfo.Login, AuthInfo.Password)
	err := row.Scan(&Res.ID)
	if err != nil {
		return nil, err
	}
	Res.Token = auth.CreateToken()
	return &Res, nil
}
func (P *PostgresRepository) GetFile(ctx context.Context, ID string) (domain.FileInfo, error) {
	var Res domain.FileInfo
	const query = `SELECT name, users, id, path, file FROM files WHERE id = $1`
	row := P.Pool.QueryRow(ctx, query, ID)
	err := row.Scan(&Res.Name, &Res.Users, &Res.ID, &Res.Path, &Res.File)
	if err != nil {
		return domain.FileInfo{}, err
	}
	return Res, nil
}

func (P *PostgresRepository) DeleteFile(ctx context.Context, ID string) error {
	const query = `DELETE FROM files WHERE id = $1;`
	_, err := P.Pool.Exec(ctx, query, ID)
	if err != nil {
		return err
	}
	return nil
}
