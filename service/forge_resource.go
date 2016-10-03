package service

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/jhaals/puppet-anvil/api"
	"github.com/pmylund/sortutil"
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

// Query modules, filter by user supplied query
func (f *ForgeResource) GetModules(w http.ResponseWriter, r *http.Request) {
	user, mod, err := parseQueryGetParam(r)
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
		Uri:      fmt.Sprintf("/v3/releases/%s/%s", metadata.Name, metadata.Version),
		Version:  metadata.Version,
		FileUri:  fmt.Sprintf("/v3/files/%s", path),
		Md5:      checksum,
		Metadata: metadata,
	}, nil
}

// Get all info about a module
func (f *ForgeResource) GetModuleInfo(w http.ResponseWriter, r *http.Request) {
	user, mod, err := parseUserModulePathParam(r)
	if err != nil {
		setBadRequestResponse(w, err)
		return
	}

	results, err := f.getModuleResult(user, mod)

	response := results

	if err := setOKResponse(w, response); err != nil {
		setInternalServerErrorResponse(w, err)
	}
}

func (f *ForgeResource) getModuleResult(user string, mod string) (api.ModuleResult, error) {
	var err error
	var result api.ModuleResult
	releases := make([]api.Release, 0)

	//get metadata for given user-modules
	modules := ListModules(filepath.Join(f.ModulePath, user, mod))
	if len(modules) == 0 {
		log.Println("list of modules is empty")
		return result, err
	}

	//sort the list of modules and get most recent version
	sortutil.DescByField(modules, "Version")
	v := modules[0].Version
	n := modules[0].Name
	m := modules[0]

	//get path and md5 of most recent version
	path := filepath.Join(user, mod, user+"-"+mod+"-"+v+".tar.gz")
	checksum, err := Checksum(filepath.Join(f.ModulePath, path))
	if err != nil {
		log.Println(err)
		return api.ModuleResult{}, fmt.Errorf("Could not get md5 of module")
	}

	//builds the releases array
	for _, metadata := range modules {
		path := filepath.Join(user, mod, user+"-"+mod+"-"+metadata.Version+".tar.gz")
		release, err := f.getReleaseFromMetaData(metadata, path)
		if err == nil {
			releases = append(releases, release)
		}
	}

	return api.ModuleResult{
		Uri:      fmt.Sprintf("/v3/modules/%s", n),
		Slug:     user + "-" + mod,
		Name:     n,
		FileUri:  fmt.Sprintf("/v3/files/%s", path),
		Version:  v,
		Md5:      checksum,
		Releases: releases,
		Owner: api.Owner{
			Uri:      fmt.Sprintf("/v3/users/%s", user),
			Slug:     user,
			Username: user,
		},
		CurrentRelease: api.CurrentRelease{
			Uri:      fmt.Sprintf("/v3/releases/%s-%s-%s", user, mod, v),
			Slug:     user + "-" + mod + "-" + v,
			Metadata: m,
			ModuleAbbreviated: api.ModuleAbbreviated{
				Uri:  fmt.Sprintf("/v3/users/%s", user),
				Slug: user + "-" + mod,
				Name: mod,
				Owner: api.Owner{
					Uri:      fmt.Sprintf("/v3/users/%s", user),
					Slug:     user,
					Username: user,
				},
			},
		},
	}, nil
}

func (f *ForgeResource) getReleaseFromMetaData(metadata api.Metadata, path string) (api.Release, error) {
	return api.Release{
		Uri:     fmt.Sprintf("/v3/releases/%s-%s", metadata.Name, metadata.Version),
		FileUri: fmt.Sprintf("/v3/files/%s", path),
		Version: metadata.Version,
		Slug:    metadata.Name + "-" + metadata.Version,
	}, nil
}

// Get information about a given release
func (f *ForgeResource) GetReleaseInfo(w http.ResponseWriter, r *http.Request) {
	user, mod, version, err := parseReleasePathParam(r)
	if err != nil {
		setBadRequestResponse(w, err)
		return
	}

	results, err := f.getReleaseResult(user, mod, version)

	response := results

	if err := setOKResponse(w, response); err != nil {
		setInternalServerErrorResponse(w, err)
	}
}

// Get information about a given release
func (f *ForgeResource) getReleaseResult(user string, mod string, version string) (api.Result, error) {
	pathToFile := filepath.Join(f.ModulePath, user, mod)
	file := user + "-" + mod + "-" + version + ".tar.gz"
	checksum, err := Checksum(filepath.Join(pathToFile, file))
	if err != nil {
		log.Println(err)
		return api.Result{}, fmt.Errorf("Could not get md5 of module")
	}

	metadata, err := GetModuleMetadata(pathToFile, file)
	if err != nil {
		log.Println(err)
		return api.Result{}, fmt.Errorf("Could not retrieve metadata for module")
	}

	return api.Result{
		Uri:      fmt.Sprintf("/v3/releases/%s-%s-%s", user, mod, version),
		Version:  version,
		FileUri:  fmt.Sprintf("/v3/files/%s/%s/%s", user, mod, file),
		Md5:      checksum,
		Metadata: metadata,
	}, nil
}
