package main

import (
	"flag"
	"fmt"
	"image-compress/utils"
	"os"
	"path/filepath"
	"sync"
)

var path string
var quality int
var outputPath string

func main() {
	flag.StringVar(&path, "p", "./org", "The path point to image")
	flag.IntVar(&quality, "q", 90, "The quality for compressing")
	flag.StringVar(&outputPath, "o", "./", "The path of output")
	flag.Parse()

	file, err := os.Open(path)
	if err != nil {
		utils.HandleError(err)
	}

	defer file.Close()

	if _, err := os.Stat(outputPath); err != nil {
		fmt.Println("Create the output path", outputPath)
		if err := os.Mkdir(outputPath, os.ModePerm); err != nil {
			utils.HandleError(err)
		}
	}

	fileInfo, err := file.Stat()
	if err != nil {
		utils.HandleError(err)
	}

	if fileInfo.IsDir() {
		files, err := os.ReadDir(path)
		if err != nil {
			utils.HandleError(err)
		}

		var wg sync.WaitGroup
		wg.Add(len(files))

		for _, f := range files {
			if f.IsDir() {
				continue
			}

			go func() {
				defer wg.Done()

				file, err := os.Open(filepath.Join(path, f.Name()))
				if err != nil {
					utils.HandleError(err)
				}

				defer file.Close()

				if err := compressFile(file); err != nil {
					utils.HandleError(err)
				}
			}()
		}

		wg.Wait()
		fmt.Println("Process complete")
	} else {
		if err := compressFile(file); err != nil {
			utils.HandleError(err)
		}
	}
}

func compressFile(file *os.File) error {
	if name, err := utils.ImageProcessing(file, quality, outputPath); err != nil {
		fmt.Println("Compress failed", name)
		return err
	}

	return nil
}
