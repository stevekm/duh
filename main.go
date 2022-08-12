package main

import (
	"code.cloudfoundry.org/bytefmt"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type SizeMapEntry struct {
	Path      string
	Size      int64
	Percent   float64
	BarLength int
	Bar       string
	ByteSize  string
	StartDir  bool
}

func NewSizeMapEntry(path string, size int64, totalSize int64, startDir string) SizeMapEntry {
	percent := CalcPercent(size, totalSize)
	barLength := CalcBarLength(percent)
	bar := CreateBar(barLength)
	byteSize := bytefmt.ByteSize(uint64(size))
	entry := SizeMapEntry{
		Path:      path,
		Size:      size,
		Percent:   percent,
		BarLength: barLength,
		Bar:       bar,
		ByteSize:  byteSize,
		StartDir:  false,
	}
	if entry.Path == startDir {
		entry.StartDir = true
	}
	return entry
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

// get the size of all items in the subdir
// https://stackoverflow.com/questions/71153302/how-to-set-depth-for-recursive-iteration-of-directories-in-filepath-walk-func
func SubDirSizes(subDirPath string) (map[string]int64, error) {
	dirSizes := map[string]int64{}
	maxDepth := 0 // do not recurse below the top level of the dir

	err := filepath.Walk(subDirPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// need to try and strip out extraneous / from the path string because we need to use the count of /'s for maxDepth count
		trimmedPath := strings.TrimLeft(path, subDirPath)
		depthCount := strings.Count(trimmedPath, string(os.PathSeparator))

		// re-trim leading / if it was present
		// trimmedPath2 := strings.TrimLeft( trimmedPath, string(os.PathSeparator))
		// depthCount := strings.Count(trimmedPath2, string(os.PathSeparator))
		// fmt.Printf("%v subDirPath: %v path: %v trimmedPath: %v\n", depthCount, subDirPath, path, trimmedPath2)

		// depthCount := strings.Count(path, string(os.PathSeparator))
		// fmt.Printf("%v %v\n", depthCount, path)

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

func FormatLines(entries []SizeMapEntry) []string {
	lines := []string{}

	// make a line for each item except start dir
	var startDirIndex int
	for i, entry := range entries {
		if entry.StartDir != true {
			var line string = entry.ByteSize + "\t" + entry.Path + "\t" + entry.Bar
			lines = append(lines, line)
		} else {
			startDirIndex = i
		}
	}

	// make start dir line
	lines = append(lines, "-----")
	var line string = entries[startDirIndex].ByteSize + "\t" + entries[startDirIndex].Path
	lines = append(lines, line)

	return lines
}

func GetDirEntries(startDir string) []SizeMapEntry {
	// remove any trailing / from the starting path
	// NOTE: This is important because otherwise it messes up the maxDepth calculation for recursion prevention
	trimmedStartDir := strings.TrimRight(startDir, string(os.PathSeparator))

	// get all the file and dir items and their sizes
	sizes, err := SubDirSizes(trimmedStartDir)
	if err != nil {
		log.Fatal(err)
	}
	totalSize := sizes[trimmedStartDir]

	sizeMapEntries := []SizeMapEntry{}
	for path, size := range sizes {
		entry := NewSizeMapEntry(path, size, totalSize, trimmedStartDir)
		sizeMapEntries = append(sizeMapEntries, entry)
	}

	// sort by path
	sort.Slice(sizeMapEntries, func(i, j int) bool {
		return sizeMapEntries[i].Path < sizeMapEntries[j].Path
	})

	return sizeMapEntries
}

func main() {
	args := os.Args[1:]
	startDir := args[0]

	sizeMapEntries := GetDirEntries(startDir)

	lines := FormatLines(sizeMapEntries)

	for _, line := range lines {
		fmt.Println(line)
	}
}
