package fileprocessing

import (
	"errors"
	"os"
	"path/filepath"
)

type FileProcessing struct {
	Dir string
}

func NewFileProc(Dir string) *FileProcessing {
	os.Mkdir(Dir, os.ModePerm)
	return &FileProcessing{Dir: Dir}
}

func (F *FileProcessing) DeleteFile(Name string) error {
	err := os.Remove(filepath.Join(F.Dir) + Name)
	return err
}

func (F *FileProcessing) CreateAndWrite(Name string, Data []byte) error {
	file, err := os.Create(filepath.Join(F.Dir) + Name)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			return err
		}
	}
	os.WriteFile(filepath.Join(F.Dir)+Name, Data, 0644)
	defer file.Close()
	return nil
}
