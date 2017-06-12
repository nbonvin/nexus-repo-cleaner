package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

var keepVersion *int

func main() {
	keepVersion = flag.Int("keep", 10, "Amount of versions to keep")
	flag.Parse()
	if *keepVersion < 1 {
		log.Fatalf("at least 1 version needs to be kept.")
		return
	}
	args := flag.Args()
	if len(args) != 1 {
		log.Fatalf("usage: repocleaner <repo_path> -keep=<amount of versions to keep>")
		return
	}
	repoPath := args[0]

	err := filepath.Walk(repoPath, visit)
	if err != nil {
		log.Fatalf("unable to walk down '%s': %v", repoPath, err)
	}
}

func visit(path string, f os.FileInfo, err error) error {
	if f.IsDir() {
		return cleanArtifactDirectory(path)
	}
	return nil
}

func cleanArtifactDirectory(dirname string) error {
	sortedDirs, err := readDirNumSort(dirname, false)
	if err != nil {
		return fmt.Errorf("unable to read folder '%s': %v", dirname, err)
	}

	versions := []os.FileInfo{}
	for _, f := range sortedDirs {

		maven, err := isMavenVersion(f)
		if err != nil {
			return fmt.Errorf("unable to get maven version: %v", err)
		}
		if maven {
			fmt.Printf("directory: %s, version: %v, last mod date: %v\n", dirname, f.Name(), f.ModTime())
			versions = append(versions, f)
		}
	}

	// skip 'keepVersion' most recent folders, keeping the other ones to delete them
	if len(versions) > *keepVersion {
		versionsToDel := versions[*keepVersion:]
		for _, f := range versionsToDel {
			path := dirname + "/" + f.Name()
			fmt.Printf("deleting directory: %s, version: %v, last mod date: %v\n", dirname, f.Name(), f.ModTime())
			err := os.RemoveAll(path)
			if err != nil {
				return fmt.Errorf("unable to remove folder '%s': %v", path, err)
			}
		}
	}
	return nil
}

func isMavenVersion(fileInfo os.FileInfo) (bool, error) {
	if !fileInfo.IsDir() {
		return false, nil
	}
	matched, err := regexp.MatchString(`^\d+\.\d+.+`, fileInfo.Name())
	if err != nil {
		return false, fmt.Errorf("unable to check maven version: %v", err)
	}
	return matched, nil

}

func readDirNumSort(dirname string, reverse bool) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	if reverse {
		sort.Sort(sort.Reverse(byDate(list)))
	} else {
		sort.Sort(byDate(list))
	}
	return list, nil
}

type byDate []os.FileInfo

func (f byDate) Len() int {
	return len(f)
}
func (f byDate) Less(i, j int) bool {
	return f[i].ModTime().Unix() > f[j].ModTime().Unix()
}
func (f byDate) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
