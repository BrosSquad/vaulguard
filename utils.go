package main

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/BrosSquad/vaulguard/log"
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

func memoryUsage(ctx context.Context, logger *log.Logger) {
	var m runtime.MemStats
	for {
		select {
		case <- ctx.Done():
			return
		case timer := <-time.After(30 * time.Second):
			runtime.ReadMemStats(&m)
			year, month, day := timer.Date()
			hour := timer.Hour()
			minute := timer.Minute()
			second := timer.Second()
			logger.Info(ctx,"---------------------------------------------")
			logger.Info(ctx,"Memory Usage")
			logger.Info(ctx,"Current Time: %d.%d.%d %d:%d:%d", day, month, year, hour, minute, second)
			logger.Info(ctx,"Current Memory usage: %v b (%v MiB)", m.Alloc, m.Alloc/1024/1024)
			logger.Info(ctx,"Total Allocations: %v b (%v MiB)", m.TotalAlloc, m.TotalAlloc/1024/1024)
			logger.Info(ctx,"Allocated from system: %v b (%v MiB)", m.Sys, m.Sys/1024/1024)
			logger.Info(ctx,"---------------------------------------------")
		}
	}
}