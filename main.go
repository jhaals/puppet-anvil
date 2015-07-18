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

	forge := &ForgeResource{ModulePath: modulePath}
	admin := &AdminResource{ModulePath: modulePath}

	mr := mux.NewRouter()
	mr.HandleFunc("/v3/files/{user}/{module}/{fileName}", forge.GetModule).Methods("GET")
	mr.HandleFunc("/v3/releases", forge.GetReleases).Methods("GET")
	mr.HandleFunc("/modules/{user}/{module}/{fileName}", admin.UpsertFile).Methods("PUT")

	http.Handle("/", mr)

	http.ListenAndServe(":"+port, nil)
}

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
	moduleName := r.URL.Query().Get("module")
	if !strings.Contains(moduleName, "-") {
		http.Error(w, "request must be /v3/releases?module=user-module", 400)
		return
	}

	user := strings.Split(moduleName, "-")[0]
	mod := strings.Split(moduleName, "-")[1]

	modules := ListModules(filepath.Join(f.ModulePath, user, mod))
	// No modules, return minimal json response.
	if len(modules) == 0 {
		fmt.Fprintf(w, `{"pagination":{"next":null},"results":[]}`)
		return
	}
	response := new(Response)
	response.Pagination = Pagination{Next: false}

	for _, metadata := range modules {
		checksum, err := Checksum(filepath.Join(f.ModulePath, user, mod, moduleName+"-"+metadata.Version+".tar.gz"))
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

type AdminResource struct {
	ModulePath string
}

func (a *AdminResource) UpsertFile(w http.ResponseWriter, r *http.Request) {
	content, _ := ioutil.ReadAll(r.Body)
	user, module, fileName, err := getModulePathComponents(r)
	if err != nil {
		rest.SetBadRequestResponse(w)
		return
	}

	fullModulePath := fmt.Sprintf("%s/%s/%s", a.ModulePath, user, module)
	if _, err := os.Stat(fullModulePath); err != nil {
		os.MkdirAll(fullModulePath, 0755)
	}

	if err := ioutil.WriteFile(fullModulePath+"/"+fileName, content, 0666); err != nil {
		log.Fatal(err)
		return
	}

	url := fmt.Sprintf("/v3/files/%s/%s/%s", user, module, fileName)

	w.Header().Set("Location", url)
	w.Header().Set("Content-Type", "application/json")
}

func getModulePathComponents(r *http.Request) (string, string, string, error) {
	user, err := rest.PathString(r, "user")
	if err != nil {
		return "", "", "", err
	}
	module, err := rest.PathString(r, "module")
	if err != nil {
		return "", "", "", err
	}
	fileName, err := rest.PathString(r, "fileName")
	if err != nil {
		return "", "", "", err
	}
	return user, module, fileName, nil
}
