package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Repo interface {
	changefeedRepo
	attachmentFinder
	attachmentUpdater
	suiteFinder
	suiteUpdater
	caseFinder
	logFinder
}

func Handler(repo Repo, publicDir string) http.Handler {
	r := mux.NewRouter()

	// middleware
	r.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	// api v1
	api := r.PathPrefix("/v1/").Subrouter()
	// attachments
	api.Path("/attachments/{id}").
		Handler(newGetAttachmentHandler(repo)).
		Methods(http.MethodGet)
	api.Path("/attachments/{id}").
		Handler(newDeleteAttachmentHandler(repo)).
		Methods(http.MethodDelete)
	api.Path("/attachments").
		Handler(newGetAttachmentCollectionHandler(repo)).
		Methods(http.MethodGet)
	api.Path("/attachments").
		Handler(newDeleteAttachmentCollectionHandler(repo)).
		Methods(http.MethodDelete)
	// suites
	api.Path("/suites/{id}").
		Handler(newGetSuiteHandler(repo)).
		Methods(http.MethodGet)
	api.Path("/suites/{id}").
		Handler(newDeleteSuiteHandler(repo)).
		Methods(http.MethodDelete)
	api.Path("/suites").
		Handler(newGetSuiteCollectionHandler(repo)).
		Methods(http.MethodGet)
	api.Path("/suites").
		Handler(newDeleteSuiteCollectionHandler(repo)).
		Methods(http.MethodDelete)
	// cases
	api.Path("/cases/{id}").
		Handler(newGetCaseHandler(repo)).
		Methods(http.MethodGet)
	api.Path("/suites/{suite_id}/cases").
		Handler(newGetCaseCollectionHandler(repo)).
		Methods(http.MethodGet)
	// logs
	api.Path("/cases/{case_id}/logs").
		Handler(newGetLogCollectionHandler(repo)).
		Methods(http.MethodGet)
	// changefeed
	api.Path("/changefeed").
		Handler(newChangefeed(repo).newHandler()).
		Methods(http.MethodGet)

	// frontend
	frontend := r.PathPrefix("/").Subrouter()
	frontend.PathPrefix("/").Handler(newFrontendHandler(publicDir))
	return r
}

func writeJson(w http.ResponseWriter, code int, msg interface{}) error {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	b, err := json.Marshal(&msg)
	if err != nil {
		log.Panicf("marshal json: %v\n", err)
	}
	if _, err := w.Write(b); err != nil {
		return fmt.Errorf("write json: %v", err)
	}
	return nil
}
