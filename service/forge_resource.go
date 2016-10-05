package service

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sort"

	"github.com/hashicorp/go-version"
	"github.com/jhaals/puppet-anvil/api"
)

// Handlers that implement the Forge API
// https://forgeapi.puppetlabs.com/
type ForgeResource struct {
	ModulePath string
}

//methods so we can use the sort.Interface for api.Release
//using hashicorp/go-version for the sorting algorithm
//perhaps there's a better way to do this, but this seems to work
type ByReleaseVersion []api.Release

func (slice ByReleaseVersion) Len() int {
	return len(slice)
}
func (slice ByReleaseVersion) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
func (slice ByReleaseVersion) Less(i, j int) bool {
	v1, _ := version.NewVersion(slice[i].Version)
	v2, _ := version.NewVersion(slice[j].Version)
	return v1.LessThan(v2)
}

//methods so we can use the sort.Interface for api.Metadata
//using hashicorp/go-version for the sorting algorithm
//perhaps there's a better way to do this, but this seems to work
type ByMetadataVersion []api.Metadata

func (slice ByMetadataVersion) Len() int {
	return len(slice)
}
func (slice ByMetadataVersion) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
func (slice ByMetadataVersion) Less(i, j int) bool {
	v1, _ := version.NewVersion(slice[i].Version)
	v2, _ := version.NewVersion(slice[j].Version)
	return v1.LessThan(v2)
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

	//gets metadata for most recent release
	current := f.sortMetaDataByVersion(modules)
	currentVersion := current[0].Version
	currentName := current[0].Name
	currentMetadata := current[0]

	//get path and md5 of most recent version
	path := filepath.Join(user, mod, user+"-"+mod+"-"+currentVersion+".tar.gz")
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

	//sorts the releases in descending order
	//even though there is a "CurrentRelease" object, it appears that some clients
	//(like librarian-puppet) use the first release in the releases object as the "Current Release"
	sr := f.sortReleasesByVersion(releases)

	return api.ModuleResult{
		Uri:      fmt.Sprintf("/v3/modules/%s", currentName),
		Slug:     user + "-" + mod,
		Name:     currentName,
		FileUri:  fmt.Sprintf("/v3/files/%s", path),
		Version:  currentVersion,
		Md5:      checksum,
		Releases: sr,
		Owner: api.Owner{
			Uri:      fmt.Sprintf("/v3/users/%s", user),
			Slug:     user,
			Username: user,
		},
		CurrentRelease: api.CurrentRelease{
			Uri:      fmt.Sprintf("/v3/releases/%s-%s-%s", user, mod, currentVersion),
			Slug:     user + "-" + mod + "-" + currentVersion,
			Metadata: currentMetadata,
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

//sorts []api.Release in descending order
func (f *ForgeResource) sortReleasesByVersion(releases []api.Release) []api.Release {
	sort.Sort(sort.Reverse(ByReleaseVersion(releases)))
	return releases
}

//sorts []api.Metadata in descending order
func (f *ForgeResource) sortMetaDataByVersion(metadata []api.Metadata) []api.Metadata {
	sort.Sort(sort.Reverse(ByMetadataVersion(metadata)))
	return metadata
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
