package fileprocessing_test

import (
	fileprocessing "midtest/internal/repo/file_processing"
	"testing"
)

func TestAll(t *testing.T) {
	Proc := fileprocessing.NewFileProc("/update")
	data := []byte("TEST")
	err := Proc.CreateAndWrite("TEST", data)
	if err != nil {
		t.Error(err.Error())
	}
	err = Proc.DeleteFile("TEST")
}
