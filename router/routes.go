package router

import (
    "net/http"
    "github.com/gorilla/mux"

    "github.com/da4nik/web-ui/api"
)

func Router() *mux.Router {
    r := mux.NewRouter()

    api_root := r.PathPrefix("/api").Subrouter()

    nodes := api_root.PathPrefix("/nodes").Subrouter()
    nodes.HandleFunc("/", api.Nodes)

    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./dist/")))

    return r
}
