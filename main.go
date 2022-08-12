package main

import (
	"fmt"
	"io/fs"
	"os"
	"log"
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
func SubDirSizes(subDirPath string) (map[string]int64, error) {
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

	return dirSizes, err
}

// func FormatSubdirStr(dirMap map[string]int64, curDir string) string {
// 	var outputStr string
// 	var totalSize = dirMap[curDir]
// 	for key, value := range dirMap {
// 		s +=
// 	}
// }


func CalcPercent(this int64, total int64) float64 {
	result := float64(this) / float64(total)
	return result
}


func main() {
	curDir := "."
	sizes, err := SubDirSizes(curDir)
	if err != nil {
			log.Fatal(err)
	}
	totalSize := sizes[curDir]
	for key, value := range sizes {
		fmt.Printf("%v %v: %v\n", key, value, CalcPercent(value, totalSize))
	}
}
