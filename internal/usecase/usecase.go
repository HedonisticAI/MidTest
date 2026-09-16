package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"midtest/internal/auth"
	"midtest/internal/domain"
	"os"
	"slices"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	ADMToken     string
	Dir          string
	PostgresRepo PostgresRepo
	CacheRepo    CacheRepo
}

func NewUsecase(ADMToken string, Dir string, PostgresRepo PostgresRepo, CacheRepo CacheRepo) Usecase {
	os.Mkdir(Dir, os.ModePerm)
	return &Service{
		Dir:          Dir,
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

func (S *Service) WriteFile(ctx context.Context, input WriteFileInput) (*WriteFileOutput, error) {
	if !S.CacheRepo.IsActive(input.Meta.Token) {
		return nil, domain.ErrUnauthorized
	}

	if input.Meta.Name == "" {
		return nil, domain.ErrBadParameter
	}

	data := input.Body
	if !input.Meta.File {
		if input.Json == nil {
			return nil, domain.ErrBadParameter
		}
		var err error
		data, err = json.Marshal(input.Json)
		if err != nil {
			return nil, fmt.Errorf("marshal json document: %w", err)
		}
	}

	if len(data) == 0 {
		return nil, domain.ErrBadParameter
	}

	if err := os.MkdirAll(S.Dir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("create storage directory: %w", err)
	}

	info := S.transform(input.Meta)
	info.ID = uuid.NewString()
	info.Path = S.Dir + string(os.PathSeparator)

	// Сначала пишем во временный файл. Так пользователь никогда не увидит
	// частично записанный документ по окончательному имени.
	tmp, err := os.CreateTemp(S.Dir, ".upload-*")
	if err != nil {
		return nil, fmt.Errorf("create temporary file: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return nil, fmt.Errorf("write document: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return nil, fmt.Errorf("close document: %w", err)
	}

	if err := S.PostgresRepo.NewFile(ctx, info); err != nil {
		return nil, err
	}

	finalName := info.Path + info.Name
	if err := os.Rename(tmpName, finalName); err != nil {
		// Если запись в БД уже создана, хотя бы не оставляем временный файл.
		// Удаление записи можно вынести в транзакцию/компенсирующую операцию
		// репозитория при дальнейшем развитии интерфейса.
		return nil, fmt.Errorf("move document to storage: %w", err)
	}

	result := &WriteFileOutput{Name: input.Meta.Name}
	if !input.Meta.File {
		result.Json = data
	}

	return result, nil
}

func (S *Service) ListFiles(ctx context.Context, List ListInput) ([]domain.FileInfo, error) {
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
		return fileTypeSwitcher(file.([]byte), Res.File), nil
	}
	memfile, err := os.ReadFile(Res.Path + Res.Name)
	if err != nil {
		return nil, err
	}
	S.CacheRepo.LoadFile(Res.Name, memfile, 10*time.Minute)
	return fileTypeSwitcher(memfile, Res.File), nil
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

func fileTypeSwitcher(file []byte, filetype bool) []byte {
	if filetype {
		return file
	} else {
		json, err := json.Marshal(file)
		if err != nil {
			return nil
		}
		return json
	}
}

func (S *Service) transform(Meta Meta) domain.FileInfo {
	return domain.FileInfo{Name: Meta.Name, File: Meta.File, Users: Meta.Users, Path: S.Dir}
}
