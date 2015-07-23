package service

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/jhaals/puppet-anvil/api"
)

// Handlers that implement the Forge API
// https://forgeapi.puppetlabs.com/
type ForgeResource struct {
	ModulePath string
}

// Serve up module archive
func (f *ForgeResource) GetModule(w http.ResponseWriter, r *http.Request) {
	user, module, fileName, err := parseFileNamePathParam(r)
	if err != nil {
		setBadRequestResponse(w, err)
		return
	}

	http.ServeFile(w, r, filepath.Join(f.ModulePath, user, module, fileName))
}

// Query releases, filter by user supplied module
func (f *ForgeResource) GetReleases(w http.ResponseWriter, r *http.Request) {
	user, mod, err := parseModuleGetParam(r)
	if err != nil {
		setBadRequestResponse(w, err)
		return
	}

	results, err := f.getResults(user, mod)

	response := &api.Response{
		Pagination: api.Pagination{
			Next: false, // nil?
		},
		Results: results,
	}

	if err := setOKResponse(w, response); err != nil {
		setInternalServerErrorResponse(w, err)
	}
}

func (f *ForgeResource) getResults(user string, mod string) ([]api.Result, error) {
	results := make([]api.Result, 0)

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

func (f *ForgeResource) getResult(metadata api.Metadata, path string) (api.Result, error) {
	checksum, err := Checksum(filepath.Join(f.ModulePath, path))
	if err != nil {
		log.Println(err)
		return api.Result{}, fmt.Errorf("not a module")
	}
	return api.Result{
		Uri:      fmt.Sprintf("/v3/release/%s/%s", metadata.Name, metadata.Version),
		Version:  metadata.Version,
		FileUri:  fmt.Sprintf("/v3/files/%s", path),
		Md5:      checksum,
		Metadata: metadata,
	}, nil
}
