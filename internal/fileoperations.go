package internal

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

var _ FileOperations = (*FileOperationsImpl)(nil)

type FileOperations interface {
	Unzip(zip []byte) (data []byte, err error)
	CreateNewTempPath(p string) (newPath string, err error)
	SaveTo(data []byte, path string) error
	MoveRunningExeToBackup(p string) error
	MoveNewExeToOriginalExe(newPath string, oldPath string) error
}

type FileOperationsImpl struct {
}

func (FileOperationsImpl) CreateNewTempPath(p string) (newPath string, err error) {
	return newPath + ".new.temp", nil
}

func (FileOperationsImpl) SaveTo(data []byte, path string) error {
	return os.WriteFile(path, data, 0755)
}

func (FileOperationsImpl) MoveRunningExeToBackup(p string) error {
	return os.Rename(p, p+".old")
}

func (FileOperationsImpl) MoveNewExeToOriginalExe(newPath string, oldPath string) error {
	return os.Rename(newPath, oldPath)
}

func (f FileOperationsImpl) GetAssetReader(url string) (data []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (f FileOperationsImpl) Unzip(data []byte) (uncompressedFile []byte, err error) {
	archive, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}

	if len(archive.File) != 1 {
		return nil, fmt.Errorf("expected 1 file in zip, got %d", len(archive.File))
	}

	file := archive.File[0]
	uzip, err := file.Open()
	if err != nil {
		return nil, err
	}
	uncompressedFile, err = io.ReadAll(uzip)
	if err != nil {
		return nil, err
	}

	return uncompressedFile, nil
}
