package result

import (
	"time"
)

type Result struct {
	FilePath string
	Name     string
	Size     int64

	Lines int
	Words int

	// Métricas de tiempo
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration

	Error error
}

func (r Result) IsSuccess() bool {
	return r.Error == nil
}

func (r Result) Summary() map[string]interface{} {
	return map[string]interface{}{
		"file":     r.Name,
		"path":     r.FilePath,
		"size":     r.Size,
		"lines":    r.Lines,
		"words":    r.Words,
		"duration": r.Duration.String(),
		"success":  r.IsSuccess(),
	}
}

func CollectStatistics(results []Result) map[string]interface{} {
	totalFiles := len(results)
	successCount := 0
	failCount := 0
	var totalSize int64
	var totalLines int
	var totalWords int
	var totalDuration time.Duration

	for _, r := range results {
		if r.IsSuccess() {
			successCount++
			totalSize += r.Size
			totalLines += r.Lines
			totalWords += r.Words
			totalDuration += r.Duration
		} else {
			failCount++
		}
	}

	var avgDuration time.Duration
	if successCount > 0 {
		avgDuration = totalDuration / time.Duration(successCount)
	}

	return map[string]interface{}{
		"total_files":      totalFiles,
		"success_count":    successCount,
		"failure_count":    failCount,
		"success_rate":     float64(successCount) / float64(totalFiles),
		"total_size":       totalSize,
		"total_lines":      totalLines,
		"total_words":      totalWords,
		"total_duration":   totalDuration.String(),
		"average_duration": avgDuration.String(),
	}
}
