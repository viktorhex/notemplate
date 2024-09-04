package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"testing"
)

func TestCreateEntry(t *testing.T) {
	// specify an enty template and a suffix
	template := "__tests"
	suffix := "examplesuffix"
	p := CreateEntryParams{template, suffix}

	// create an entry
	create_entry(p)

	// the created dir is the same name as the template used
	dir := template

	// expect first items to be dirs
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, folder := range files {
		fmt.Println(folder.Name())
		if folder.IsDir() != true {
			t.Fatalf("expected dir")
		}
		// then remove all created test files
		removeAll(t, dir, folder)
	}
	// optionally assert folder/file names and structure
}

// dir is directory of entries
// file is an entry, iow file is a directory of toml files
func removeAll(t *testing.T, dir string, file fs.DirEntry) {
	innerFiles, err := os.ReadDir(path.Join(dir, file.Name()))
	if err != nil {
		log.Fatal(err)
	}
	for _, innerFile := range innerFiles {
		err := os.Remove(path.Join(dir, file.Name(), innerFile.Name()))
		if err != nil {
			t.Fatalf("error removing test files %v", err.Error())
		}
	}
	err = os.Remove(path.Join(dir, file.Name()))
	if err != nil {
		t.Fatalf("error removing test folders %v", err.Error())
	}
}
