package router

import (
    "github.com/ant0ine/go-json-rest/rest"
    "log"
    "net/http"

    "github.com/da4nik/web-ui/api"
)

var (
    HttpLogin    = "admin"
    HttpPassword = "admin"
)

func getRestApi() *rest.Api {
    restapi := rest.NewApi()
    restapi.Use(rest.DefaultDevStack...)

    router, err := rest.MakeRouter(
        rest.Get("/nodes", api.GetAllNodes),
        rest.Get("/nodes/:node/services", api.GetNodeServices),

        rest.Get("/services", api.GetAllServices),

        rest.Put("/services/new", api.PutServiceConfig),
        rest.Get("/services/:service/config", api.GetServiceConfig),
        rest.Post("/services/:service/config", api.PostServiceConfig),
        rest.Delete("/services/:service/config", api.DeleteServiceConfig),

        rest.Post("/services/:service/build", api.PostBuildService),     // GET param nodes
        rest.Post("/services/:service/restart", api.PostRestartService), //GET param node
    )
    if err != nil {
        log.Fatal(err)
    }
    restapi.SetApp(router)
    return restapi
}

func Setup() {
    restapi := getRestApi()
    restapi.Use(&rest.AuthBasicMiddleware{
        Realm: "test zone",
        Authenticator: func(userId string, password string) bool {
            if userId == HttpLogin && password == HttpPassword {
                return true
            }
            return false
        },
    })

    http.Handle("/api/", http.StripPrefix("/api", restapi.MakeHandler()))

    http.Handle("/", http.FileServer(http.Dir("./dist/")))
}
