package utils

import (
	"context"
	"runtime"
	"time"

	"github.com/BrosSquad/vaulguard/log"
)

func MemoryUsage(ctx context.Context, sleep time.Duration, logger *log.Logger) {
	var m runtime.MemStats
	for {
		select {
		case <-ctx.Done():
			return
		case timer := <-time.After(sleep):
			runtime.ReadMemStats(&m)
			year, month, day := timer.Date()
			hour := timer.Hour()
			minute := timer.Minute()
			second := timer.Second()
			logger.Info(ctx, "---------------------------------------------")
			logger.Info(ctx, "Memory Usage")
			logger.Info(ctx, "Current Time: %d.%d.%d %d:%d:%d", day, month, year, hour, minute, second)
			logger.Info(ctx, "Current Memory usage: %v b (%v MiB)", m.Alloc, m.Alloc/1024/1024)
			logger.Info(ctx, "Total Allocations: %v b (%v MiB)", m.TotalAlloc, m.TotalAlloc/1024/1024)
			logger.Info(ctx, "Allocated from system: %v b (%v MiB)", m.Sys, m.Sys/1024/1024)
			logger.Info(ctx, "---------------------------------------------")
		}
	}
}
