package api

import (
    "net/http"
)

func Nodes(w http.ResponseWriter, req *http.Request) {
    w.Write([]byte("Supernodes list"))
}
