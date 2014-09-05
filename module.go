package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
func ListModules(path string) []Metadata {
	var result []Metadata
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println(err)
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".tar.gz") {
			err := ExtractMetadata(file, path)
			if err != nil {
				log.Println(err)
				continue
			}
			metadata, err := readMetadata(filepath.Join(path, file.Name()+".metadata"))
			if err != nil {
				log.Println(err)
				continue
			}
			result = append(result, metadata)
		}
	}
	return result
}

//Extract metadata from module
func ExtractMetadata(module os.FileInfo, path string) error {
	moduleFile := filepath.Join(path, module.Name())
	metadataPath := filepath.Join(path, module.Name()+".metadata")
	metadataFile, err := os.Stat(metadataPath)

	if err == nil {
		if metadataFile.ModTime().After(module.ModTime()) {
			// Fresh metadata, skipping
			return nil
		}
	}
	log.Println("Extracting metadata.json from", moduleFile)
	fi, err := os.Open(moduleFile)
	if err != nil {
		return err
	}
	defer fi.Close()
	fz, err := gzip.NewReader(fi)
	if err != nil {
		return err
	}
	defer fz.Close()

	// tar.gz data
	s, err := ioutil.ReadAll(fz)
	if err != nil {
		return err
	}
	// TODO Prettify this thing...
	r := bytes.NewReader(s)
	tr := tar.NewReader(r)

	// Iterate through the files in the archive.
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return errors.New("metadata.json not found in " + moduleFile)
		}
		if err != nil {
			return err
		}
		// Found metadata.json, no need to read any further.
		if hdr.Name == strings.TrimRight(module.Name(), "tar.gz")+"/metadata.json" {
			f, err := os.Create(metadataPath)
			defer f.Close()
			if err != nil {
				return err
			}
			io.Copy(f, tr)
			return nil
		}
	}
}

func readMetadata(file string) (Metadata, error) {
	var m Metadata
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return m, err
	}
	json.Unmarshal(data, &m)
	return m, nil
}
