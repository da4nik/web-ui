package api

import (
    "github.com/ant0ine/go-json-rest/rest"
    porterConsul "github.com/da4nik/porter/consul"
    "net/http"
    "sync"
)

var lock = sync.RWMutex{}

func isError(err error, w rest.ResponseWriter) bool {
    if _, ok := err.(porterConsul.NoConfigError); ok {
        rest.Error(w, "Service config does not exist", http.StatusNotFound)
        return true
    }
    if err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return true
    }
    return false
}
