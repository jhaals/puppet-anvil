package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func parseFileNamePathParam(r *http.Request) (string, string, string, error) {
	fileName := mux.Vars(r)["fileName"]
	user, mod, err := parseFileName(fileName)
	return user, mod, fileName, err
}

func parseUserModulePathParam(r *http.Request) (string, string, error) {
	userModule := mux.Vars(r)["user-module"]
	user, mod, err := parseModuleName(userModule)
	return user, mod, err
}

func parseModuleGetParam(r *http.Request) (string, string, error) {
	moduleName := r.URL.Query().Get("module")
	if len(moduleName) == 0 {
		return "", "", fmt.Errorf("GET param 'module' empty")
	}
	return parseModuleName(moduleName)
}

func parseQueryGetParam(r *http.Request) (string, string, error) {
	queryName := r.URL.Query().Get("query")
	if len(queryName) == 0 {
		return "", "", fmt.Errorf("GET param 'query' empty")
	}
	return parseModuleName(queryName)
}
func parseFileName(fileName string) (string, string, error) {
	if !strings.HasSuffix(fileName, ".tar.gz") {
		return "", "", fmt.Errorf("module '%s' should be of the form '{user}-{module}-{version}.tar.gz'", fileName)
	}
	if strings.Count(fileName, "-") != 2 {
		return "", "", fmt.Errorf("module '%s' should be of the form '{user}-{module}-{version}.tar.gz'", fileName)
	}

	moduleName := fileName[0:strings.LastIndex(fileName, "-")]

	return parseModuleName(moduleName)
}

func parseModuleName(moduleName string) (string, string, error) {
	if strings.Count(moduleName, "-") != 1 {
		return "", "", fmt.Errorf("module '%s' should be of the form '{user}-{module'}", moduleName)
	}
	user := strings.Split(moduleName, "-")[0]
	mod := strings.Split(moduleName, "-")[1]

	return user, mod, nil
}

func parseReleasePathParam(r *http.Request) (string, string, string, error) {
	userModuleVersion := mux.Vars(r)["user-module-version"]
	if strings.Count(userModuleVersion, "-") != 2 {
		return "", "", "", fmt.Errorf("user-module-version '%s' should be of the form '{user}-{module}-{version}", userModuleVersion)
	}
	user := strings.Split(userModuleVersion, "-")[0]
	mod := strings.Split(userModuleVersion, "-")[1]
	version := strings.Split(userModuleVersion, "-")[2]

	return user, mod, version, nil
}
