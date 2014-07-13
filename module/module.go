package module

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
)

type Dependencies struct {
	Name         string `json:"name"`
	VersionRange string `json:"version_requirement"`
}

type Metadata struct {
	Name         string         `json:"name"`
	Version      string         `json:"version"`
	Author       string         `json:"author"`
	Licence      string         `json:"license"`
	Dependencies []Dependencies `json:"dependencies"`
}

type Result struct {
	Uri      string `json:"uri"`
	FileUri  string `json:"file_uri"`
	Version  string `json:"version"`
	Md5      string `json:"file_md5"`
	Metadata `json:"metadata"`
}
type Pagination struct {
	Next bool `json:"next"`
}
type Response struct {
	Results    []Result `json:"results"`
	Pagination `json:"pagination"`
}

// Checksum file
func Checksum(file string) string {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", md5.Sum(data))
}

// ListModules returns all tar.gz files
func ListModules(path string) []string {
	var result []string
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		match, _ := regexp.MatchString(".tar.gz$", file.Name())
		if match {
			result = append(result, filepath.Join(path, file.Name()))
			ExtractMetadata(file, path)
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(result)))
	return result
}

//Extract metadata from module
func ExtractMetadata(module os.FileInfo, path string) {
	filePath := filepath.Join(path, module.Name())
	metadata_path := filepath.Join(path, module.Name()+".metadata")

	metadataFile, err := os.Stat(metadata_path)
	if err == nil {
		if metadataFile.ModTime().After(module.ModTime()) {
			// Fresh metadata, skipping
			return
		}
	}
	// Must be GNU tar
	// TODO: Use built in gzip library.
	metadata, err := exec.Command("tar", "-z", "-x", "--wildcards", "-O", "-f", filePath, "*/metadata.json").Output()
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(metadata_path, []byte(metadata), 0644)
}

func ReadMetadata(file string) (Metadata, error) {
	var m Metadata
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return m, errors.New("Failed to read " + file)
	}
	json.Unmarshal(data, &m)
	return m, nil
}
