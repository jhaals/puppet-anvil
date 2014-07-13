package main

import (
	"encoding/json"
	"fmt"
	"github.com/jhaals/go-puppet-forge/module"
	"net/http"
	"path/filepath"
	"strings"
)

func main() {
	// List all releases for a module
	http.HandleFunc("/v3/files/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/var/lib/go-puppet-forge/modules/"+r.URL.Path[10:])
	})
	http.HandleFunc("/v3/releases", ReleaseHandler)
	http.ListenAndServe(":8080", nil)
}

func ReleaseHandler(w http.ResponseWriter, r *http.Request) {
	moduleDir := "/var/lib/go-puppet-forge/modules/"
	moduleName := r.URL.Query().Get("module")
	if moduleName == "" {
		http.Error(w, "request must be /v3/releases?module=user-module", 400)
		return
	}

	user := strings.Split(moduleName, "-")[0]
	mod := strings.Split(moduleName, "-")[1]

	modules := module.ListModules(filepath.Join(moduleDir, user, mod))
	// No modules, return minimal json response.
	if len(modules) == 0 {
		fmt.Fprintf(w, `{"pagination":{"next":null},"results":[]}`)
		return
	}
	response := new(module.Response)
	response.Pagination = module.Pagination{Next: false}

	for _, file := range modules {
		metadata, _ := module.ReadMetadata(file + ".metadata")
		var result = module.Result{
			Uri:     fmt.Sprintf("/v3/release/%s/%s", metadata.Name, metadata.Version),
			Version: metadata.Version,
			FileUri: fmt.Sprintf("/v3/files/%s/%s/%s-%s.tar.gz", user, mod, moduleName, metadata.Version),
			Md5:     module.Checksum(file)}
		result.Metadata = metadata
		response.Results = append(response.Results, result)
	}
	jsonData, _ := json.Marshal(response)
	fmt.Fprintf(w, string(jsonData))
}
