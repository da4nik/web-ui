package main

import (
    "net/http"

    "github.com/da4nik/web-ui/router"
)

func main() {
    panic( http.ListenAndServe(":8000", router.Router()) )
}
