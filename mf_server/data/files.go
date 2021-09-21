package Data

import (
	"io/fs"
	"path/filepath"
)

// Find files with a extension and returns absolut path
func Find(root, ext string) []string {
	var files []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil { return e }
		if filepath.Ext(d.Name()) == ext {
			files = append(files, root + "/" +d.Name())
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

