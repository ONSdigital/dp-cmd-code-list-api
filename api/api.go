package api

import (
	"encoding/json"
	"net/http"

	"github.com/ONSdigital/dp-code-list-api/datastore"
	"github.com/ONSdigital/go-ns/log"
	"github.com/gorilla/mux"
)

// CodeListAPI holds all endpoints which are used to access the code list resources
type CodeListAPI struct {
	router *mux.Router
	store  datastore.DataStore
}

// CreateCodeListAPI returns a constructed code list api
func CreateCodeListAPI(route *mux.Router, store datastore.DataStore) *CodeListAPI {
	api := CodeListAPI{router: route, store: store}
	api.router.HandleFunc("/code-lists", api.getCodeLists).Methods("GET")
	api.router.HandleFunc("/code-lists/{id}", api.getCodeList).Methods("GET")
	api.router.HandleFunc("/code-lists/{id}/codes", api.getCodes).Methods("GET")
	api.router.HandleFunc("/code-lists/{id}/codes/{code}", api.getCode).Methods("GET")
	return &api
}

func (c *CodeListAPI) getCodeLists(w http.ResponseWriter, r *http.Request) {

	codeLists, err := c.store.GetCodeLists()
	if err != nil {
		handleError(w, err)
		return
	}

	codeLists.Count = len(codeLists.Items)
	codeLists.Limit = len(codeLists.Items)

	bytes, err := json.Marshal(codeLists)
	if err != nil {
		handleError(w, err)
		return
	}

	writeBody(w, bytes)
	log.Debug("get all code lists", nil)
}

func (c *CodeListAPI) getCodeList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	codeList, err := c.store.GetCodeList(id)
	if err != nil {
		handleError(w, err)
		return
	}

	bytes, err := json.Marshal(codeList)
	if err != nil {
		handleError(w, err)
		return
	}

	writeBody(w, bytes)
	log.Debug("get a code list", log.Data{"id": id})
}

func (c *CodeListAPI) getCodes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	codes, err := c.store.GetCodes(id)
	if err != nil {
		handleError(w, err)
		return
	}

	codes.Count = len(codes.Items)
	codes.Limit = len(codes.Items)

	bytes, err := json.Marshal(codes)
	if err != nil {
		handleError(w, err)
		return
	}

	writeBody(w, bytes)
	log.Debug("get all codes", log.Data{"id": id})
}

func (c *CodeListAPI) getCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	codeID := vars["code"]
	code, err := c.store.GetCode(id, codeID)
	if err != nil {
		handleError(w, err)
		return
	}

	bytes, err := json.Marshal(code)
	if err != nil {
		handleError(w, err)
		return
	}

	writeBody(w, bytes)
	log.Debug("get a code", log.Data{"id": id, "code": codeID})
}

func handleError(w http.ResponseWriter, err error) {
	log.Error(err, nil)
	if datastore.NOT_FOUND == err {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func writeBody(w http.ResponseWriter, bytes []byte) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(bytes); err != nil {
		log.Error(err, nil)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
