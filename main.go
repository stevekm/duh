package main

import (
	"fmt"
	"io/fs"
	"os"
	// "log"
	"path/filepath"
	"strings"
)

// get the size of one dir
// https://stackoverflow.com/questions/32482673/how-to-get-directory-total-size
func DirSize(dirPath string) (int64, error) {
	var size int64
	err := filepath.Walk(dirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

// size, err := DirSize(startDir)
// if err != nil {
// 		log.Fatal(err)
// }

// get the size of all subdirs
// https://stackoverflow.com/questions/71153302/how-to-set-depth-for-recursive-iteration-of-directories-in-filepath-walk-func
func SubDirSizes(subDirPath string) error {
	dirSizes := map[string]int64{}
	maxDepth := 0
	err := filepath.Walk(subDirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && strings.Count(path, string(os.PathSeparator)) > maxDepth {
			// fmt.Println("skip", path)
			return fs.SkipDir
		}
		if info.IsDir() {
			// size += info.Size()
			subDirSize, _ := DirSize(path) // TODO: handle err ...
			dirSizes[path] = subDirSize

		}
		return err
	})

	fmt.Printf("dirSizes: %v\n", dirSizes)

	return err
}

func main() {
	startDir := "."
	SubDirSizes(startDir)
}
