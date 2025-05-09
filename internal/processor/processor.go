package processor

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/MrXetwy21/worker-pool/pkg/result"
)

type FileTask struct {
	FilePath string
	Verbose  bool
}

func (ft FileTask) Process() (result.Result, error) {
	res := result.Result{
		FilePath:  ft.FilePath,
		StartTime: time.Now(),
	}

	file, err := os.Open(ft.FilePath)
	if err != nil {
		return res, fmt.Errorf("error al abrir archivo: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return res, fmt.Errorf("error al obtener información del archivo: %w", err)
	}
	res.Size = fileInfo.Size()
	res.Name = filepath.Base(ft.FilePath)

	scanner := bufio.NewScanner(file)
	lineCount := 0
	wordCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++
		words := strings.Fields(line)
		wordCount += len(words)

		if ft.Verbose && lineCount%1000 == 0 {
			fmt.Printf("Procesando %s: %d líneas\n", res.Name, lineCount)
		}
	}

	if err := scanner.Err(); err != nil {
		return res, fmt.Errorf("error al escanear archivo: %w", err)
	}

	res.Lines = lineCount
	res.Words = wordCount
	res.EndTime = time.Now()
	res.Duration = res.EndTime.Sub(res.StartTime)

	return res, nil
}

func EstimateWorkload(dirPath string) (int64, error) {
	var totalSize int64

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("error al estimar carga de trabajo: %w", err)
	}

	return totalSize, nil
}

func OptimizeWorkerCount(fileCount int) int {
	numCPU := 1

	if fileCount < numCPU {
		return fileCount
	} else if fileCount < numCPU*2 {
		return numCPU
	} else {
		return numCPU + 1
	}
}
