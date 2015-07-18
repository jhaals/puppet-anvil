package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/benschw/opin-go/rest"
)

type ForgeResource struct {
	ModulePath string
}

func (f *ForgeResource) GetModule(w http.ResponseWriter, r *http.Request) {
	user, module, fileName, err := getModulePathComponents(r)
	if err != nil {
		rest.SetBadRequestResponse(w)
		return
	}

	http.ServeFile(w, r, filepath.Join(f.ModulePath, user, module, fileName))
}

func (f *ForgeResource) GetReleases(w http.ResponseWriter, r *http.Request) {
	user, mod, err := parseModuleGetParam(r)
	if err != nil {
		http.Error(w, "request must be /v3/releases?module=user-module", 400)
		return
	}

	results, err := f.getResults(user, mod)

	response := &Response{
		Pagination: Pagination{
			Next: false, // nil?
		},
		Results: results,
	}

	if err := rest.SetOKResponse(w, response); err != nil {
		rest.SetInternalServerErrorResponse(w, err)
	}
}

func (f *ForgeResource) getResults(user string, mod string) ([]Result, error) {
	results := make([]Result, 0)

	modules := ListModules(filepath.Join(f.ModulePath, user, mod))

	for _, metadata := range modules {
		path := filepath.Join(user, mod, user+"-"+mod+"-"+metadata.Version+".tar.gz")
		result, err := f.getResult(metadata, path)
		if err == nil {
			results = append(results, result)
		}
	}
	return results, nil
}

func (f *ForgeResource) getResult(metadata Metadata, path string) (Result, error) {
	checksum, err := Checksum(filepath.Join(f.ModulePath, path))
	if err != nil {

		log.Println(err)
		return Result{}, fmt.Errorf("not a module")
	}
	return Result{
		Uri:      fmt.Sprintf("/v3/release/%s/%s", metadata.Name, metadata.Version),
		Version:  metadata.Version,
		FileUri:  fmt.Sprintf("/v3/files/%s", path),
		Md5:      checksum,
		Metadata: metadata,
	}, nil
}

func parseModuleGetParam(r *http.Request) (string, string, error) {
	moduleName := r.URL.Query().Get("module")
	if !strings.Contains(moduleName, "-") {
		return "", "", fmt.Errorf("bad get param")
	}

	user := strings.Split(moduleName, "-")[0]
	mod := strings.Split(moduleName, "-")[1]

	return user, mod, nil
}
