package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sort"
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
	maxDepth := 1 // do not recurse below the top level of the dir 

	err := filepath.Walk(subDirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// need to try and strip out extraneous / from the path string because we need to use the count of /'s for maxDepth count
		trimmedPath := strings.TrimLeft(path, subDirPath)
		depthCount := strings.Count(trimmedPath, string(os.PathSeparator))
		// fmt.Printf("%v %v %v %v\n", depthCount, subDirPath, path, trimmedPath)
		
		if depthCount > maxDepth { // info.IsDir() && depthCount > maxDepth // fmt.Println("skip", path)
			return fs.SkipDir
		}

		if info.IsDir() {
			subDirSize, err := DirSize(path)
			if err != nil {
				log.Fatal(err)
			}
			dirSizes[path] = subDirSize

		} else {
			dirSizes[path] = info.Size()
		}
		return err
	})

	// fmt.Printf("dirSizes: %v\n", dirSizes)
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


func FormatLines(entries []SizeMapEntry) []string {
	lines := []string{}
	for _, entry := range entries {
		var line string = entry.Path + "\t" + entry.Bar 
		lines = append(lines, line)
	}

	// sort by path before pre-pending the byte size
	sort.Strings(lines)

	sortedLines := []string{}
	for i, line := range lines {
		var line = entries[i].ByteSize + "\t" + line
		sortedLines = append(sortedLines, line)
	}
	
	return sortedLines
}

func main() {
	args := os.Args[1:]
	startDir := args[0] 

	// remove any trailing / from the starting path 
	// NOTE: This is important because otherwise it messes up the maxDepth calculation for recursion prevention
	trimmedStartDir := strings.TrimRight(startDir, string(os.PathSeparator))

	sizes, err := SubDirSizes(trimmedStartDir)
	if err != nil {
		log.Fatal(err)
	}
	totalSize := sizes[trimmedStartDir]

	sizeMapEntries := FormatMap(sizes, totalSize)

	lines := FormatLines(sizeMapEntries)

	for _, line := range lines {
		fmt.Println(line)
	}
}
