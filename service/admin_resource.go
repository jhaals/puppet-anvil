package service

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type AdminResource struct {
	ModulePath string
}

func (a *AdminResource) UpsertFile(w http.ResponseWriter, r *http.Request) {
	content, _ := ioutil.ReadAll(r.Body)
	user, module, fileName, err := parseFileNamePathParam(r)

	if err != nil {
		SetBadRequestResponse(w, err)
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
