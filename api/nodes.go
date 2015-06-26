package api

import (
    "github.com/ant0ine/go-json-rest/rest"
)

func GetAllNodes(w rest.ResponseWriter, r *rest.Request) {
    nodes, err := Nodes()
    if isError(err, w) {
        return
    }
    w.WriteJson(&nodes)
}
