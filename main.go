package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"code.cloudfoundry.org/bytefmt"
)

type SizeMapEntry struct {
	Path      string
	Size      int64
	Percent   float64
	BarLength int
	Bar       string
	ByteSize string
}

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
			subDirSize, err := DirSize(path) // TODO: handle err ...
			if err != nil {
				log.Fatal(err)
			}
			dirSizes[path] = subDirSize

		}
		return err
	})

	fmt.Printf("dirSizes: %v\n", dirSizes)
	return dirSizes, err
}

// value between 0 and 1
func CalcPercent(this int64, total int64) float64 {
	result := float64(this) / float64(total)
	return result
}

// bar length should be between 1 and 100
// TODO: should 80 be the max length?
func CalcBarLength(percent float64) int {
	result := int(percent * 100.0)
	if result < 1 {
		result = 1
	}
	return result
}

func CreateBar(length int) string {
	result := strings.Repeat("|", length)
	return result
}

func FormatMap(sizes map[string]int64, totalSize int64) []SizeMapEntry {
	sizeMapEntries := []SizeMapEntry{}
	for key, value := range sizes {
		percent := CalcPercent(value, totalSize)
		barLength := CalcBarLength(percent)
		bar := CreateBar(barLength)
		byteSize := bytefmt.ByteSize(uint64(value))
		entry := SizeMapEntry{
			Path:      key,
			Size:      value,
			Percent:   percent,
			BarLength: barLength,
			Bar:       bar,
			ByteSize: byteSize,
		}
		sizeMapEntries = append(sizeMapEntries, entry)
	}
	return sizeMapEntries
}

func main() {
	args := os.Args[1:]
	startDir := args[0] // curDir := "."
	sizes, err := SubDirSizes(startDir)
	if err != nil {
		log.Fatal(err)
	}
	totalSize := sizes[startDir]

	sizeMapEntries := FormatMap(sizes, totalSize)

	for _, entry := range sizeMapEntries {
		fmt.Printf("%v\n", entry)
	}
}
