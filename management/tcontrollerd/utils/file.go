package utils

import (
	"bytes"
	"fmt"
	"io"
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

// RenameWriteFile writes data to filename by writing a temporary file next to
// it and renaming that file into place.
//
// os.CreateTemp creates the file with O_CREATE|O_EXCL and mode 0600, so a file
// or symlink planted at the temporary path beforehand cannot be followed, and
// the rename replaces a symlink that was left at filename instead of writing
// through it. The previous implementation guessed the name with math/rand and
// wrote through ioutil.WriteFile, which follows symlinks.
func RenameWriteFile(filename string, data []byte, perm os.FileMode) error {
	return RenameCopyFromIO(bytes.NewReader(data), filename, perm)
}

// RenameCopyFromIO streams src into a temporary file created next to dstPath
// and renames it over dstPath once the whole content has been written.
//
// This is the write primitive of the package: opening dstPath in place, as
// CopyFileFromIO used to do with O_CREATE|O_WRONLY|O_TRUNC and no O_NOFOLLOW,
// lets anyone who can create an entry in the destination directory redirect the
// write to another file and truncate it. tcontrollerd runs as root and writes
// the nginx configuration of every site into a path derived from the site id,
// so the destination names are predictable and the target of such a link would
// be overwritten with root privileges.
//
// The rename also makes the update atomic: a reader either sees the complete
// old file or the complete new one, never a half written mixture.
func RenameCopyFromIO(src io.Reader, dstPath string, perm os.FileMode) error {
	tmpFile, err := os.CreateTemp(filepath.Dir(dstPath), filepath.Base(dstPath)+".tmp.")
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
	if _, err = io.Copy(tmpFile, src); err != nil {
		tmpFile.Close()
		return err
	}
	if err = tmpFile.Close(); err != nil {
		return err
	}
	if err = os.Rename(tmpName, dstPath); err != nil {
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

// EnsureWriteFile writes a file with the given mode, creating its parent
// directory when needed.
//
// It goes through RenameCopyFromIO, so the write is atomic and cannot be
// redirected through a symlink: ioutil.WriteFile, which this used to call,
// truncates the target of a pre-existing link at the destination path.
func EnsureWriteFile(path string, data []byte, mode os.FileMode) error {
	return EnsureRenameWriteFile(path, data, mode)
}

func CopyFile(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer func(srcFile *os.File) {
		err := srcFile.Close()
		if err != nil {

		}
	}(srcFile)

	fileInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	return CopyFileFromIO(srcFile, dstPath, fileInfo.Mode())
}

// CopyFileFromIO copies src into dstPath.
//
// The copy is written to a temporary file and renamed into place, so it is
// atomic and cannot be redirected through a symlink left at dstPath. Opening
// dstPath with O_CREATE|O_WRONLY|O_TRUNC and no O_NOFOLLOW, which this used to
// do, truncated whatever a planted link pointed at.
func CopyFileFromIO(src io.Reader, dstPath string, perm os.FileMode) error {
	if err := EnsureFileDir(dstPath); err != nil {
		return err
	}

	return RenameCopyFromIO(src, dstPath, perm)
}

func CopyFileIfNotExist(srcPath, dstPath string) error {
	if exist, err := FileExist(dstPath); err != nil {
		return err
	} else if !exist {
		return CopyFile(srcPath, dstPath)
	} else {
		return nil
	}
}
