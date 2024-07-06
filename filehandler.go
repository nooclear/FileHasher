package FileHasher

import (
	"os"
)

var directories []string
var files []string

func getEntries(path string) {
	if entries, err := os.ReadDir(path); err != nil {
		panic(err)
	} else {
		for _, en := range entries {
			if en.IsDir() {
				directories = append(directories, path+"/"+en.Name())
				getEntries(path + "/" + en.Name())
			} else {
				files = append(files, path+"/"+en.Name())
			}
		}
	}
}
