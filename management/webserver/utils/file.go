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
	// os.CreateTemp creates the file with O_CREATE|O_EXCL and mode 0600, so a
	// file or symlink that was planted at the temporary path beforehand cannot
	// be followed and the payload is never readable by other users before the
	// rename. The previous implementation built a predictable name with
	// math/rand and wrote through ioutil.WriteFile, which does both.
	tmpFile, err := os.CreateTemp(filepath.Dir(filename), filepath.Base(filename)+".tmp.")
	if err != nil {
		return err
	}
	tmpName := tmpFile.Name()

	renamed := false
	defer func() {
		if !renamed {
			os.Remove(tmpName)
		}
	}()

	if err = tmpFile.Chmod(perm); err != nil {
		tmpFile.Close()
		return err
	}
	if _, err = tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return err
	}
	if err = tmpFile.Close(); err != nil {
		return err
	}
	if err = os.Rename(tmpName, filename); err != nil {
		return err
	}

	renamed = true
	return nil
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
