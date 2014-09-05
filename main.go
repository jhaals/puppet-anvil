package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	port := os.Getenv("PORT")
	modulePath := os.Getenv("MODULEPATH")

	if len(port) == 0 {
		log.Fatal("Missing PORT environment variable")
	}
	if len(modulePath) == 0 {
		log.Fatal("Missing MODULEPATH environment variable")
	}

	log.Println("Starting go-puppet-forge on port", port, "serving modules from", modulePath)

	http.HandleFunc("/v3/files/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(modulePath, r.URL.Path[10:]))
	})
	http.HandleFunc("/v3/releases", func(w http.ResponseWriter, r *http.Request) {
		ReleaseHandler(w, r, modulePath)
	})
	http.ListenAndServe(":"+port, nil)
}

func ReleaseHandler(w http.ResponseWriter, r *http.Request, modulePath string) {
	moduleName := r.URL.Query().Get("module")
	if !strings.Contains(moduleName, "-") {
		http.Error(w, "request must be /v3/releases?module=user-module", 400)
		return
	}

	user := strings.Split(moduleName, "-")[0]
	mod := strings.Split(moduleName, "-")[1]

	modules := ListModules(filepath.Join(modulePath, user, mod))
	// No modules, return minimal json response.
	if len(modules) == 0 {
		fmt.Fprintf(w, `{"pagination":{"next":null},"results":[]}`)
		return
	}
	response := new(Response)
	response.Pagination = Pagination{Next: false}

	for _, metadata := range modules {
		checksum, err := Checksum(filepath.Join(modulePath, user, mod, moduleName+"-"+metadata.Version+".tar.gz"))
		if err != nil {
			// Unable to checksum modulefile, log and skip.
			log.Println(err)
			continue
		}
		var result = Result{
			Uri:     fmt.Sprintf("/v3/release/%s/%s", metadata.Name, metadata.Version),
			Version: metadata.Version,
			FileUri: fmt.Sprintf("/v3/files/%s/%s/%s-%s.tar.gz", user, mod, moduleName, metadata.Version),
			Md5:     checksum}
		result.Metadata = metadata
		response.Results = append(response.Results, result)
	}
	jsonData, _ := json.Marshal(response)
	fmt.Fprintf(w, string(jsonData))
}
