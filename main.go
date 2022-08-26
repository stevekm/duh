package main

import (
	"code.cloudfoundry.org/bytefmt"
	"fmt"
	"github.com/TwiN/go-color"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"sort"
	"strings"
)

var logger = log.New(os.Stderr, "", 0)

type SizeMapEntry struct {
	Path      string  // original file path
	Size      int64   // original size in bytes
	Percent   float64 // percent of total dir size that this entry takes up (should be value between 0-1)
	BarLength int     // how long of a text graphic to draw
	Bar       string  // text graphic for the entry
	StartDir  bool    // if this item was the starting directory for search
}

func NewSizeMapEntry(path string, size int64, totalSize int64, startDir string) SizeMapEntry {
	percent := CalcPercent(size, totalSize)
	barLength := CalcBarLength(percent)
	bar := CreateBar(barLength)
	entry := SizeMapEntry{
		Path:      path,
		Size:      size,
		Percent:   percent,
		BarLength: barLength,
		Bar:       bar,
		StartDir:  false,
	}
	if entry.Path == startDir {
		entry.StartDir = true
	}
	return entry
}

// get the size of all items in the subdir
// https://stackoverflow.com/questions/71153302/how-to-set-depth-for-recursive-iteration-of-directories-in-filepath-walk-func
func SubDirSizes(subDirPath string) (map[string]int64, error) {
	dirSizes := map[string]int64{}
	subDirPathParts := strings.Split(subDirPath, string(os.PathSeparator))
	subDirPathPartsLen := len(subDirPathParts)

	// make sure the map value for the root dir is initialized
	dirSizes[subDirPath] += int64(0)

	err := filepath.WalkDir(subDirPath, func(path string, dirEntry fs.DirEntry, err error) error {
		// path = 001/subdir.1/file1.1163766069
		parts := strings.Split(path, string(os.PathSeparator)) // [001 subdir.1 file1.1163766069]
		root := parts[0]                                       // 001
		// fmt.Printf("subDirPath: %v, path: %v, parts: %v, root: %v\n", subDirPath, path, parts, root)

		// make sure the map value for the root dir is initialized
		dirSizes[root] += int64(0)

		// skip item that cannot be read
		if os.IsPermission(err) {
			logger.Printf("Skipping path that could not be read %q: %v\n", path, err)
			return filepath.SkipDir
		}

		// return other errors encountered
		if err != nil {
			return err
		}

		if !dirEntry.IsDir() {
			info, err := dirEntry.Info()
			if err != nil {
				return err
			}

			// dir1 as input subDirPath
			if root == subDirPath {
				key := parts[subDirPathPartsLen:][0]
				dirSizes[key] += info.Size()
			} else if subDirPathPartsLen > 1 {
				// dir1/go as input subDirPath
				key := parts[subDirPathPartsLen:][0]
				dirSizes[key] += info.Size()
			} else {
				// . as input subDirPath
				dirSizes[root] += info.Size()
			}

			dirSizes[subDirPath] += info.Size()
		}
		return err
	})

	// somehow a stray '' is getting in the map during test cases, remove it
	// TODO: figure out why this is happening, probably due to indexing ^^^ up there or something...
	if _, ok := dirSizes[""]; ok {
		delete(dirSizes, "")
	}

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

// make the text graphic that will be displayed
func CreateBar(length int) string {
	result := strings.Repeat("|", length)
	return result
}

func GetByteSizeColor(size int64) string {
	var col string
	if size >= 1024*1024*1024*1024 { // T
		col = color.Purple
	} else if size >= 1024*1024*1024 { // G
		col = color.Cyan
	} else if size >= 1024*1024 { // M
		col = color.Red
	} else if size >= 1024 { // K
		col = color.Yellow
	} else {
		col = color.Gray
	}
	return col
}

func GetPercentColor(percent float64) string {
	var col string
	if percent >= 0.80 {
		col = color.Purple
	} else if percent >= 0.60 {
		col = color.Cyan
	} else if percent >= 0.40 {
		col = color.Red
	} else if percent >= 0.20 {
		col = color.Yellow
	} else {
		col = color.Gray
	}
	return col
}

func FormatBar(bar string, percent float64) string {
	col := GetPercentColor(percent)
	barStr := color.Ize(col, bar)
	return barStr
}

func FormatSize(size int64) string {
	// sizeStr := color.Ize(color.Bold, bytefmt.ByteSize(uint64(size)))
	sizeStr := bytefmt.ByteSize(uint64(size))
	col := GetByteSizeColor(size)
	sizeStr = color.Ize(col, sizeStr)
	return sizeStr
}

func FormatEntryLine(entry SizeMapEntry) string {
	sizeStr := FormatSize(entry.Size)
	var line string = sizeStr + "\t" + entry.Path + "\t" + FormatBar(entry.Bar, entry.Percent)
	return line
}

func FormatStartDirLine(entry SizeMapEntry) string {
	line := color.Ize(color.Bold, bytefmt.ByteSize(uint64(entry.Size))) + "\t" + entry.Path
	return line
}

// build all the lines of text that should be printed to the console
func FormatLines(entries []SizeMapEntry) []string {
	lines := []string{}

	// make a line for each item except start dir
	var startDirIndex int
	for i, entry := range entries {
		if entry.StartDir != true {
			// NOTE: consider printing just the basename instead of the full path // path.Base(entry.Path)
			line := FormatEntryLine(entry)
			lines = append(lines, line)
		} else {
			startDirIndex = i
		}
	}

	// make start dir line
	lines = append(lines, "-----")
	var line string = FormatStartDirLine(entries[startDirIndex])
	lines = append(lines, line)

	return lines
}

func GetDirEntries(startDir string) []SizeMapEntry {
	// remove any trailing / from the starting path
	// NOTE: This is important because otherwise it messes up the maxDepth calculation for recursion prevention
	cleanPath := path.Clean(startDir)

	// get all the file and dir items and their sizes
	sizes, err := SubDirSizes(cleanPath)
	if err != nil {
		log.Fatal(err)
	}
	totalSize := sizes[cleanPath]

	sizeMapEntries := []SizeMapEntry{}
	for path, size := range sizes {
		entry := NewSizeMapEntry(path, size, totalSize, cleanPath)
		sizeMapEntries = append(sizeMapEntries, entry)
	}

	// sort by path
	sort.Slice(sizeMapEntries, func(i, j int) bool {
		return sizeMapEntries[i].Path < sizeMapEntries[j].Path
	})

	return sizeMapEntries
}

// https://pkg.go.dev/runtime/pprof
// https://github.com/google/pprof/blob/main/doc/README.md
// $ go tool pprof cpu.prof
// $ go tool pprof mem.prof
// (pprof) top
func StartProfiler() (*os.File, *os.File) {
	cpuFile, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	// defer cpuFile.Close() // error handling omitted for example
	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	// defer pprof.StopCPUProfile()

	memFile, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	// defer memFile.Close() // error handling omitted for example
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(memFile); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}

	return cpuFile, memFile
}

func main() {
	var enableProfile bool = false
	args := os.Args[1:]
	var startDir string
	if len(args) < 1 {
		startDir = "."
	} else {
		startDir = args[0]
	}
	if len(args) > 1 {
		enableProfile = true
	}
	if enableProfile {
		cpuFile, memFile := StartProfiler()
		defer cpuFile.Close()
		defer memFile.Close()
		defer pprof.StopCPUProfile()
	}

	sizeMapEntries := GetDirEntries(startDir)

	lines := FormatLines(sizeMapEntries)

	for _, line := range lines {
		fmt.Println(line)
	}
}
