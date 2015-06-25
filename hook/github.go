package hook

import (
    "github.com/ant0ine/go-json-rest/rest"
    consulapi "github.com/da4nik/porter/consul"
    "net/http"
)

func PostPushHook(w rest.ResponseWriter, r *rest.Request) {
    var payload map[string]interface{}
    err := r.DecodeJsonPayload(&payload)
    if err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    repoUrl := payload["repository"].(map[string]interface{})["clone_url"].(string)
    lastCommit := payload["head_commit"].(map[string]interface{})["id"].(string)
    configs, err := consulapi.ListConfigs()
    if err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    for _, config := range configs {
        if config.Repo == repoUrl {
            config.LastCommit = lastCommit
            config.Update()
            w.WriteHeader(http.StatusOK)
            return
        }
    }
    rest.Error(w, "Can't find service to update", http.StatusNotFound)
}
