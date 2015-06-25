package api

import (
    "github.com/ant0ine/go-json-rest/rest"
    porterConsul "github.com/da4nik/porter/consul"
    "net/http"
)

func GetServiceConfig(w rest.ResponseWriter, r *rest.Request) {
    serviceName := r.PathParam("service")
    lock.RLock()
    defer lock.RUnlock()
    config, err := porterConsul.GetServiceConfig(serviceName)
    if isError(err, w) {
        return
    }
    w.WriteJson(&config)
}

func PutServiceConfig(w rest.ResponseWriter, r *rest.Request) {
    config := new(porterConsul.ServiceConfig)
    err := r.DecodeJsonPayload(config)
    if isError(err, w) {
        return
    }
    if config.Repo == "" {
        rest.Error(w, "repo url required", 400)
        return
    }
    if config.Name == "" {
        rest.Error(w, "service name required", 400)
        return
    }
    lock.Lock()
    defer lock.Unlock()
    err = config.Update()
    if isError(err, w) {
        return
    }
    w.WriteJson(&config)
}

func PostServiceConfig(w rest.ResponseWriter, r *rest.Request) {
    serviceName := r.PathParam("service")
    config, err := porterConsul.GetServiceConfig(serviceName)
    if isError(err, w) {
        return
    }
    err = r.DecodeJsonPayload(config)
    if isError(err, w) {
        return
    }
    if config.Repo == "" {
        rest.Error(w, "repo url required", 400)
        return
    }
    config.Name = serviceName
    lock.Lock()
    defer lock.Unlock()
    err = config.Update()
    if err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteJson(&config)
}

func DeleteServiceConfig(w rest.ResponseWriter, r *rest.Request) {
    serviceName := r.PathParam("service")
    config, err := porterConsul.GetServiceConfig(serviceName)
    if isError(err, w) {
        return
    }
    lock.Lock()
    defer lock.Unlock()
    err = config.Delete()
    if isError(err, w) {
        return
    }
    w.WriteHeader(http.StatusOK)
}
