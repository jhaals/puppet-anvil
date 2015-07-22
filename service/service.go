package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/benschw/puppet-anvil/api"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Create a new AnvilService
func New(port string, modulePath string) *AnvilService {
	return &AnvilService{
		Bind:       ":" + port,
		ModulePath: modulePath,
	}
}

// Manage http server providing private Puppet Forge functionality
type AnvilService struct {
	Bind       string
	ModulePath string
}

// Wire routing for http server and start it up
func (s *AnvilService) Run() error {
	log.Println("Starting Puppet Anvil")

	forge := &ForgeResource{ModulePath: s.ModulePath}
	admin := &AdminResource{ModulePath: s.ModulePath}

	mr := mux.NewRouter()
	mr.HandleFunc("/v3/files/{user}/{module}/{fileName}", forge.GetModule).Methods("GET")
	mr.HandleFunc("/v3/releases", forge.GetReleases).Methods("GET")
	mr.HandleFunc("/admin/module/{fileName}", admin.UpsertFile).Methods("PUT")

	http.Handle("/", handlers.LoggingHandler(os.Stdout, mr))

	return http.ListenAndServe(s.Bind, nil)
}

func setOKResponse(w http.ResponseWriter, entity interface{}) error {
	return setResponse(w, entity, http.StatusOK)
}
func setBadRequestResponse(w http.ResponseWriter, e error) {
	resp := api.NewErrorResponse(e)
	err := setResponse(w, resp, http.StatusBadRequest)
	if err != nil {
		log.Print(err)
		setResponse(w, nil, http.StatusBadRequest)
	}
}
func setInternalServerErrorResponse(w http.ResponseWriter, e error) {
	resp := api.NewErrorResponse(e)
	err := setResponse(w, resp, http.StatusInternalServerError)
	if err != nil {
		log.Print(err)
		setResponse(w, nil, http.StatusBadRequest)
	}
}

func setResponse(w http.ResponseWriter, entity interface{}, code int) error {
	var body string
	if entity != nil {
		b, err := json.Marshal(entity)
		if err != nil {
			return err
		}
		body = string(b[:])
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if entity != nil {
		fmt.Fprint(w, body)
	}
	return nil
}
