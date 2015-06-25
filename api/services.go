package api

import (
    "github.com/ant0ine/go-json-rest/rest"
    "net/http"
)

func GetAllServices(w rest.ResponseWriter, r *rest.Request) {
    lock.RLock()
    defer lock.RUnlock()
    services, err := Services()
    if isError(err, w) {
        return
    }
    w.WriteJson(&services)
}

func GetNodeServices(w rest.ResponseWriter, r *rest.Request) {
    nodeName := r.PathParam("node")
    lock.RLock()
    defer lock.RUnlock()
    services, err := NodeServices(nodeName)
    if isError(err, w) {
        return
    }
    w.WriteJson(&services)
}

func PostBuildService(w rest.ResponseWriter, r *rest.Request) {
    serviceName := r.PathParam("service")
    nodes := r.URL.Query()["nodes"]
    lock.Lock()
    defer lock.Unlock()
    err := FireBuildEvent(serviceName, nodes)
    if isError(err, w) {
        return
    }
    w.WriteHeader(http.StatusOK)
}

func PostRestartService(w rest.ResponseWriter, r *rest.Request) {
    serviceName := r.PathParam("service")
    nodes := r.URL.Query()["nodes"]
    lock.Lock()
    defer lock.Unlock()
    err := FireRestartEvent(serviceName, nodes)
    if isError(err, w) {
        return
    }
    w.WriteHeader(http.StatusOK)
}
