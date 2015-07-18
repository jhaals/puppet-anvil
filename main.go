package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/benschw/opin-go/rest"
	"github.com/gorilla/mux"
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

	log.Println("Starting Puppet Anvil on port", port, "serving modules from", modulePath)

	mr := mux.NewRouter()
	mr.HandleFunc("/v3/files/{user}/{module}/{fileName}", func(w http.ResponseWriter, r *http.Request) {
		log.Println("files")

		user, _ := rest.PathString(r, "user")
		module, _ := rest.PathString(r, "module")
		fileName, _ := rest.PathString(r, "fileName")

		http.ServeFile(w, r, filepath.Join(modulePath, user, module, fileName))
	}).Methods("GET")

	mr.HandleFunc("/v3/releases", func(w http.ResponseWriter, r *http.Request) {
		log.Println("releases")
		ReleaseHandler(w, r, modulePath)
	}).Methods("GET")

	mr.HandleFunc("/modules/{user}/{module}/{fileName}", func(w http.ResponseWriter, r *http.Request) {
		log.Println("upload")
		upsertFile(w, r)
	}).Methods("PUT")

	http.Handle("/", mr)

	http.ListenAndServe(":"+port, nil)
}
func upsertFile(w http.ResponseWriter, r *http.Request) {
	content, _ := ioutil.ReadAll(r.Body)
	user, _ := rest.PathString(r, "user")
	module, _ := rest.PathString(r, "module")
	fileName, _ := rest.PathString(r, "fileName")

	modulePath := fmt.Sprintf("/modules/%s/%s", user, module)
	if _, err := os.Stat(modulePath); err != nil {
		os.MkdirAll(modulePath, 0755)
	}

	if err := ioutil.WriteFile(modulePath+"/"+fileName, content, 0666); err != nil {
		log.Fatal(err)
		return
	}

	url := fmt.Sprintf("/v3/files/%s/%s/%s", user, module, fileName)

	w.Header().Set("Location", url)
	w.Header().Set("Content-Type", "application/json")
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
