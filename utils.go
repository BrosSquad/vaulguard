package main

import (
	"os"
	"path"
	"path/filepath"
)

func getAbsolutePath(file string) (string, error) {
	if !path.IsAbs(file) {
		out, err := filepath.Abs(file)
		if err != nil {
			return "", err
		}
		return out, nil
	}

	return file, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
func dirExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func createDirs(paths ...string) error {
	for _, p := range paths {
		dir := filepath.Dir(p)
		if !dirExists(dir) {
			if err := os.MkdirAll(dir, DefaultPermission); err != nil {
				return err
			}
		}
	}

	return nil
}
