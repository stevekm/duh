package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"github.com/google/go-cmp/cmp"
	"testing"
)

// 
// HELPER FUNCTIONS FOR TESTS
// 
// create a temp file in a dir and write something to its contents
func createTempFile(tempdir string, filename string, contents string) (*os.File, string) {
	tempfile, err := os.CreateTemp(tempdir, filename)
	if err != nil {
		log.Fatal(err)
	}
	// defer tempfile.Close()

	// write to the file
	if contents != "" {
		nbytesWritten, err := tempfile.WriteString(contents)
		if err != nil {
			fmt.Println(nbytesWritten)
			log.Fatal(err)
		}

		// need to reset the cursor after writing
		i, err := tempfile.Seek(0, 0)
		if err != nil {
			fmt.Println("Error", err, i)
			log.Fatal(err)
		}
	}

	// get the randomly generated file basename
	fi, err := tempfile.Stat()
	if err != nil {
		log.Fatal(err)
	}
	basename := fi.Name()

	return tempfile, basename
}

func createSubDir(tempdir string, filename string) string {
	subdir := filepath.Join(tempdir, filename)
	err := os.MkdirAll(subdir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	return subdir
}


// set up a bunch of temp files and subdirs to use in test cases
func createTempFilesDirs1(tempdir string) ([]string, []*os.File) {
	subdir1 := createSubDir(tempdir, "subdir.1")
	subdir2 := createSubDir(tempdir, "subdir.2")
	subdir3 := createSubDir(tempdir, "subdir.3")

	tempfile1, _ := createTempFile(subdir1, "file1.", "writes\n")
	// defer tempfile1.Close()

	tempfile2, _ := createTempFile(subdir2, "file2.", ".........")
	// defer tempfile2.Close()

	// create this file in the root dir
	tempfile3, tempfile3Basename := createTempFile(tempdir, "file3.", "foobarfoobar")
	// defer tempfile3.Close()

	// duplicate file
	tempfile4, _ := createTempFile(subdir2, tempfile3Basename, "sometextgoeshere")
	// defer tempfile4.Close()

	tempfile5, _ := createTempFile(subdir3, "file5.", "blahblahblahblahblah")
	// defer tempfile5.Close()

	tempDirs := []string{subdir1, subdir2, subdir3}
	tempFiles := []*os.File{tempfile1, tempfile2, tempfile3, tempfile4, tempfile5}

	return tempDirs, tempFiles
}





// 
// TEST CASES
// 
// test that the find method only detects the expected files
func TestFindAllFiles(t *testing.T) {
	// automatically gets cleaned up when all tests end
	tempdir := t.TempDir() 
	// looks like this;
	// /var/folders/y4/8rsn2mvj5qv2mk5v5gyk8d2c0000gq/T/TestFindAllFiles3634983381/001
	tempDirs, tempFiles := createTempFilesDirs1(tempdir)

	tests := map[string]struct {
		input string
		want map[string]int64
	}{
		"first": {
			input: tempdir, 
			want: map[string]int64{
				tempdir: 64,
				tempFiles[2].Name(): 12, // file3.
				tempDirs[0]: 7, // subdir.1
				tempDirs[1]: 25,// subdir.2
				tempDirs[2]: 20, // subdir.3
			},
		},
		"second": {
			input: tempdir +  string(os.PathSeparator), 
			want: map[string]int64{
				tempdir: 64,
				tempFiles[2].Name(): 12, // file3.
				tempDirs[0]: 7, // subdir.1
				tempDirs[1]: 25,// subdir.2
				tempDirs[2]: 20, // subdir.3
			},
		},
		
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T){
			sizes, err := SubDirSizes(tempdir)
			if err != nil {
				log.Fatal(err)
			}
			got := sizes
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("got vs want mismatch (-want +got):\n%s", diff)
			}
		})
	}

}