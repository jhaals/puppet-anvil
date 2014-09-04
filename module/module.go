package module

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

// ListModules returns all tar.gz files
func ListModules(path string) []string {
	var result []string
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println(err)
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".tar.gz") {
			result = append(result, filepath.Join(path, file.Name()))
			ExtractMetadata(file, path)
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(result)))
	return result
}

//Extract metadata from module
func ExtractMetadata(module os.FileInfo, path string) {
	moduleFile := filepath.Join(path, module.Name())
	metadataPath := filepath.Join(path, module.Name()+".metadata")
	metadataFile, err := os.Stat(metadataPath)

	if err == nil {
		if metadataFile.ModTime().After(module.ModTime()) {
			// Fresh metadata, skipping
			return
		}
	}
	log.Println("Extracting metadata.json from", moduleFile)
	fi, err := os.Open(moduleFile)
	if err != nil {
		log.Println(err)
		return
	}
	defer fi.Close()

	fz, err := gzip.NewReader(fi)
	if err != nil {
		log.Println(err)
		return
	}
	defer fz.Close()

	// tar.gz data
	s, err := ioutil.ReadAll(fz)
	if err != nil {
		log.Println(err)
		return
	}
	// TODO Prettify this thing...
	r := bytes.NewReader(s)
	tr := tar.NewReader(r)

	// Iterate through the files in the archive.
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			log.Println(err)
			break
		}
		// Found metadata.json, no need to read any further.
		if hdr.Name == strings.TrimRight(module.Name(), "tar.gz")+"/metadata.json" {
			f, err := os.Create(metadataPath)
			defer f.Close()
			if err != nil {
				log.Println(err)
			}
			io.Copy(f, tr)
			break
		}
	}
}

func ReadMetadata(file string) Metadata {
	var m Metadata
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
		return m
	}
	json.Unmarshal(data, &m)
	return m
}
