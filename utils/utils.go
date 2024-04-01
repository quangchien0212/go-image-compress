package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/h2non/bimg"
	"github.com/h2non/filetype"
)

func ImageProcessing(file *os.File, quality int, dirname string) (string, error) {
	fileInfo, err := file.Stat()
	buffer := make([]byte, fileInfo.Size())
	file.Read(buffer)

	if isImage, err := IsImage(buffer); err != nil || !isImage {
		return fileInfo.Name(), err
	}

	fmt.Println("Compressing", file.Name())

	filename := fileNameWithoutExtension(fileInfo.Name()) + ".webp"

	converted, err := bimg.NewImage(buffer).Convert(bimg.WEBP)
	if err != nil {
		return fileInfo.Name(), err
	}

	processed, err := bimg.NewImage(converted).Process(bimg.Options{Quality: quality})
	if err != nil {
		return filename, err
	}

	writeError := bimg.Write(fmt.Sprintf(dirname+"/%s", filename), processed)
	if writeError != nil {
		return filename, writeError
	}

	return filename, nil
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func IsImage(buf []byte) (bool, error) {
	mime, err := filetype.Match(buf)
	if err != nil {
		return false, err
	}

	if mime.MIME.Type == "image" {
		return true, nil
	}

	return false, nil
}

func HandleError(err error) (b bool) {
	if err != nil {
		// notice that we're using 1, so it will actually log where
		// the error happened, 0 = this function, we don't want that.
		_, filename, line, _ := runtime.Caller(1)
		log.Fatalf("[error] %s:%d %v", filename, line, err)
		b = true
	}
	return
}
