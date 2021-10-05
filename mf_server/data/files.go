package Data

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
)

func CreateFile(absolutePath string) *os.File{
	file, err := os.Create(absolutePath)
	if err != nil {
		fmt.Println("Error creating file: \n"+err.Error())
		return nil
	}
	_ = file.Close()
	return file
}

// Check if file exists
func FileExists(absolutePath string) bool {
	_, err := os.Stat(absolutePath)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}

// Find files with a extension and returns absolut path
func FindFilesByExtension(rootPath, ext string) []string {
	var files []string
	filepath.WalkDir(rootPath, func(s string, d fs.DirEntry, e error) error {
		if e != nil { return e }
		if filepath.Ext(d.Name()) == ext {
			files = append(files, s)
		}
		return nil
	})
	return files
}

func RetrieveFileInformation(rootPath, ext string) []FileInformation{
	var fileInfos []FileInformation
	var files = FindFilesByExtension(rootPath, ext)
	for _, absolutePath := range files {
		// get last modified time
		file, err := os.Stat(absolutePath)
		if err != nil { break }
		fileInfo := FileInformation{
			AbsolutePath: absolutePath,
			Timestamp:    file.ModTime().Format("02-01-2006 15:04:05"),
			Size: strconv.Itoa(int(file.Size())),
		}
		fileInfos = append(fileInfos, fileInfo)
	}
	return fileInfos
}

func ReplaceExt(filename, repExt string) string {
	var extension = filepath.Ext(filename)
	var name = filename[0:len(filename)-len(extension)] + "." + repExt
	return name
}

