package Data

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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

func GetPathFromString(absolutePath string) string{
	var paths  = strings.Split(absolutePath, "/")
	paths = paths[:len(paths)-1]
	var finalpath = strings.Join(paths, "/")
	return finalpath
}

func GetFilenameFromString(absoulutePath string) string{
	var pathname  = strings.Split(absoulutePath, "/")
	var filename = pathname[len(pathname)-1]
	return filename
}

func ReplaceExt(filename, repExt string) string {
	var extension = filepath.Ext(filename)
	var name = filename[0:len(filename)-len(extension)] + "." + repExt
	return name
}

