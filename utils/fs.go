package utils

import (
	"os"
	"path"
	"path/filepath"
)

func GetAbsolutePath(file string) (string, error) {
	if !path.IsAbs(file) {
		out, err := filepath.Abs(file)
		if err != nil {
			return "", err
		}
		return out, nil
	}

	return file, nil
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func DirExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func CreateDirs(permission os.FileMode, paths ...string) error {
	for _, p := range paths {
		if !DirExists(p) {
			if err := os.MkdirAll(p, permission); err != nil {
				return err
			}
		}
	}

	return nil
}
