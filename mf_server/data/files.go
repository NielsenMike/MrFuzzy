package Data

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func CreateFile(absolutePath string) *os.File{
	file, err := os.Create(absolutePath)
	if err != nil {
		fmt.Println("Error creating file: \n"+err.Error())
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
			files = append(files, rootPath + "/" +d.Name())
		}
		return nil
	})
	return files
}

func ReplaceExt(filename, repExt string) string {
	var extension = filepath.Ext(filename)
	var name = filename[0:len(filename)-len(extension)] + "." + repExt
	return name
}

