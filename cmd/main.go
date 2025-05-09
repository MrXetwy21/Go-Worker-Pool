package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/MrXetwy21/worker-pool/internal/pool"
	"github.com/MrXetwy21/worker-pool/internal/processor"
)

func main() {
	// Configuración mediante flags
	dirFlag := flag.String("dir", "./data", "Directorio con archivos para procesar")
	workersFlag := flag.Int("workers", 5, "Número de workers concurrentes")
	verboseFlag := flag.Bool("verbose", false, "Mostrar información detallada")
	flag.Parse()

	if _, err := os.Stat(*dirFlag); os.IsNotExist(err) {
		log.Fatalf("El directorio %s no existe", *dirFlag)
	}

	var filePaths []string
	err := filepath.Walk(*dirFlag, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			filePaths = append(filePaths, path)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error al recorrer el directorio: %v", err)
	}

	if len(filePaths) == 0 {
		log.Printf("No se encontraron archivos para procesar en %s", *dirFlag)
		return
	}

	log.Printf("Iniciando procesamiento de %d archivos con %d workers", len(filePaths), *workersFlag)

	p := pool.New(*workersFlag)
	startTime := time.Now()

	tasks := make([]processor.FileTask, len(filePaths))
	for i, path := range filePaths {
		tasks[i] = processor.FileTask{
			FilePath: path,
			Verbose:  *verboseFlag,
		}
	}

	p.Process(tasks)

	results := p.Wait()

	duration := time.Since(startTime)
	successCount := 0
	errorCount := 0

	for _, res := range results {
		if res.Error != nil {
			errorCount++
			if *verboseFlag {
				log.Printf("Error procesando %s: %v", res.FilePath, res.Error)
			}
		} else {
			successCount++
			if *verboseFlag {
				log.Printf("Procesado correctamente: %s (líneas: %d, tamaño: %d bytes)",
					res.FilePath, res.Lines, res.Size)
			}
		}
	}

	fmt.Println("\n--- Resumen de procesamiento ---")
	fmt.Printf("Archivos procesados exitosamente: %d\n", successCount)
	fmt.Printf("Archivos con errores: %d\n", errorCount)
	fmt.Printf("Tiempo total: %v\n", duration)
	fmt.Printf("Promedio por archivo: %v\n", duration/time.Duration(len(filePaths)))
}
