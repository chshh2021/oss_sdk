package lib

import (
	"fmt"
	"time"
)

/**
 * path:bucket/group/year/month/day/microunixtime_rand(1000).jpg
 *
 */
func GenerationPath(group, suffix string) string {
	now := time.Now()
	year, month, day := now.Date()
	microunixtime := now.UnixNano()
	path := fmt.Sprintf("/%s/%d_%d_%d_%v.%s", group, year, month, day, microunixtime, suffix)
	return path
}

// 2025-03-18 Tiger
func GetSuffixByMimeType(suffix string) map[string]string {
	typeInfo := make(map[string]string)
	for k, v := range MIMETYPE {
		typeInfo[v] = k
	}
	return typeInfo
}
