package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func EnsureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.FileMode(0755))
	}
	return nil
}

func EnsureFileDir(path string) error {
	return EnsureDir(filepath.Dir(path))
}

func FileExist(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	} else {
		if stat.IsDir() {
			return false, fmt.Errorf("%s is dir", path)
		} else {
			return true, nil
		}
	}
}

func FilesExist(paths ...string) (bool, error) {
	for _, path := range paths {
		exist, err := FileExist(path)
		if err != nil {
			return false, err
		}
		if !exist {
			return false, nil
		}
	}
	return true, nil
}

func RenameWriteFile(filename string, data []byte, perm os.FileMode) error {
	randFileName := filename + ".tmp." + RandStr(8)
	if err := ioutil.WriteFile(randFileName, data, perm); err != nil {
		return err
	}
	return os.Rename(randFileName, filename)
}

func EnsureRenameWriteFile(path string, data []byte, mode os.FileMode) error {
	err := EnsureFileDir(path)
	if err != nil {
		return err
	}

	return RenameWriteFile(path, data, mode)
}

func EnsureWriteFile(path string, data []byte, mode os.FileMode) error {
	err := EnsureFileDir(path)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, mode)
}
