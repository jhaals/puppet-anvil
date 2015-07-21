package service

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/benschw/puppet-anvil/api"
)

// Handlers to help manage private forge
type AdminResource struct {
	ModulePath string
}

// PUT a module archive
func (a *AdminResource) UpsertFile(w http.ResponseWriter, r *http.Request) {
	content, _ := ioutil.ReadAll(r.Body)
	user, module, fileName, err := parseFileNamePathParam(r)

	if err != nil {
		setBadRequestResponse(w, err)
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

	modResp := &api.AdminModule{
		FileUri: fmt.Sprintf("/v3/files/%s/%s/%s", user, module, fileName),
	}

	if err := setOKResponse(w, modResp); err != nil {
		setInternalServerErrorResponse(w, err)
	}
}
