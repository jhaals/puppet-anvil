package main

import (
	"log"
	"net/http"

	"github.com/benschw/opin-go/ophttp"
	"github.com/benschw/opin-go/rest"
	"github.com/gorilla/mux"
)

func NewAnvilService(port string, modulePath string) *AnvilService {
	return &AnvilService{
		Server:     ophttp.NewServer(":" + port),
		ModulePath: modulePath,
	}
}

type AnvilService struct {
	Server     *ophttp.Server
	ModulePath string
}

func (s *AnvilService) Run() error {
	log.Println("Starting Puppet Anvil on port")

	forge := &ForgeResource{ModulePath: s.ModulePath}
	admin := &AdminResource{ModulePath: s.ModulePath}

	mr := mux.NewRouter()
	mr.HandleFunc("/v3/files/{user}/{module}/{fileName}", forge.GetModule).Methods("GET")
	mr.HandleFunc("/v3/releases", forge.GetReleases).Methods("GET")
	mr.HandleFunc("/admin/{user}/{module}/{fileName}", admin.UpsertFile).Methods("PUT")

	http.Handle("/", mr)

	err := s.Server.Start()
	return err
}

func (s *AnvilService) Stop() {
	log.Println("Stopping Puppet Anvil")
	s.Server.Stop()
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
