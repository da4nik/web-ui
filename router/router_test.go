package router

import (
    "fmt"
    "github.com/ant0ine/go-json-rest/rest/test"
    "github.com/da4nik/web-ui/api"
    "github.com/stretchr/testify/assert"
    "io/ioutil"
    "net/http"
    "testing"
)

const (
    testServiceName = "test_service"
)

type nodes []struct {
    Node    string
    Address string
}

func getNodes(t *testing.T) nodes {
    var result nodes
    recorded := test.RunRequest(t, http.StripPrefix("/api", getRestApi().MakeHandler()),
        test.MakeSimpleRequest("GET", "http://localhost/api/nodes", nil))
    recorded.CodeIs(200)
    recorded.ContentTypeIsJson()
    err := recorded.DecodeJsonPayload(&result)
    if err != nil {
        t.Fatal(err)
    }
    return result
}

func TestListNodes(t *testing.T) {
    getNodes(t)
}

func TestListServices(t *testing.T) {
    assert := assert.New(t)
    recorded := test.RunRequest(t, http.StripPrefix("/api", getRestApi().MakeHandler()),
        test.MakeSimpleRequest("GET", "http://localhost/api/services", nil))
    recorded.CodeIs(200)
    recorded.ContentTypeIsJson()
    var result map[string]interface{}
    err := recorded.DecodeJsonPayload(&result)
    if err != nil {
        t.Fatal(err)
    }
    _, ok := result["consul"]
    assert.False(ok)
}

func TestListNodeServices(t *testing.T) {
    nodes, err := api.Nodes()
    if err != nil {
        t.Fatal(err)
    }
    if len(nodes) < 1 {
        t.Fatal("Need at least one node to be registered in consul")
    }
    nodeName := nodes[0].Node
    recorded := test.RunRequest(t, http.StripPrefix("/api", getRestApi().MakeHandler()),
        test.MakeSimpleRequest("GET", fmt.Sprintf("http://localhost/api/nodes/%s/services", nodeName), nil))
    recorded.CodeIs(200)
    recorded.ContentTypeIsJson()
}

func TestGetServiceConfig(t *testing.T) {
    assert := assert.New(t)
    services, err := api.Services()
    if err != nil {
        t.Fatal(err)
    }
    if len(services) < 1 {
        t.Fatal("Need at least one service config to be stored in consul kv")
    }
    var serviceName string
    for serviceName, _ = range services {
        break
    }
    recorded := test.RunRequest(t, http.StripPrefix("/api", getRestApi().MakeHandler()),
        test.MakeSimpleRequest("GET", fmt.Sprintf("http://localhost/api/services/%s/config", serviceName), nil))
    recorded.CodeIs(200)
    recorded.ContentTypeIsJson()

    var result map[string]interface{}
    err = recorded.DecodeJsonPayload(&result)
    if err != nil {
        t.Fatal(err)
    }
    _, ok := result["Value"]
    assert.False(ok)
}

func getConfig(serviceName string, t *testing.T) (result map[string]interface{}) {
    recorded := test.RunRequest(t, http.StripPrefix("/api", getRestApi().MakeHandler()),
        test.MakeSimpleRequest("GET", fmt.Sprintf("http://localhost/api/services/%s/config", serviceName), nil))
    recorded.CodeIs(200)
    recorded.ContentTypeIsJson()
    err := recorded.DecodeJsonPayload(&result)
    if err != nil {
        t.Fatal(err)
    }
    return
}

func putServiceConfig(config map[string]interface{}, t *testing.T) {
    recorded := test.RunRequest(t, http.StripPrefix("/api", getRestApi().MakeHandler()),
        test.MakeSimpleRequest("PUT", "http://localhost/api/services/new", config))
    recorded.CodeIs(200)
    recorded.ContentTypeIsJson()
}

func TestPutServiceConfig(t *testing.T) {
    assert := assert.New(t)
    config := map[string]interface{}{
        "Volumes": []string{"/a:/b"},
        "Ports":   []string{"64:58"},
        "Env":     []string{"TEST=yes"},
        "Name":    testServiceName,
        "Repo":    "http://site.ru/a.git",
    }
    putServiceConfig(config, t)

    storedConfig := getConfig(testServiceName, t)
    assert.Equal(config["Name"], storedConfig["Name"])
    assert.Equal(config["Repo"], storedConfig["Repo"])
}

func TestPostServiceConfig(t *testing.T) {
    assert := assert.New(t)
    storedConfig := getConfig(testServiceName, t)
    storedConfig["LastCommit"] = "commitId-asdsdfsdf"
    recorded := test.RunRequest(t, http.StripPrefix("/api", getRestApi().MakeHandler()),
        test.MakeSimpleRequest("POST", fmt.Sprintf("http://localhost/api/services/%s/config", storedConfig["Name"].(string)), storedConfig))

    recorded.CodeIs(200)
    recorded.ContentTypeIsJson()
    result := getConfig(testServiceName, t)
    assert.Equal("commitId-asdsdfsdf", result["LastCommit"].(string))
}

func TestDeleteServiceConfig(t *testing.T) {
    recorded := test.RunRequest(t, http.StripPrefix("/api", getRestApi().MakeHandler()),
        test.MakeSimpleRequest("GET", fmt.Sprintf("http://localhost/api/services/%s/config", testServiceName), nil))
    recorded.CodeIs(200)

    recorded = test.RunRequest(t, http.StripPrefix("/api", getRestApi().MakeHandler()),
        test.MakeSimpleRequest("DELETE", fmt.Sprintf("http://localhost/api/services/%s/config", testServiceName), nil))

    recorded.CodeIs(200)

    recorded = test.RunRequest(t, http.StripPrefix("/api", getRestApi().MakeHandler()),
        test.MakeSimpleRequest("GET", fmt.Sprintf("http://localhost/api/services/%s/config", testServiceName), nil))
    recorded.CodeIs(404)
}

func TestSendBuildEvent(t *testing.T) {
    config := map[string]interface{}{
        "Volumes":    []string{},
        "Ports":      []string{"9000:8000"},
        "Env":        []string{"TEST=yes"},
        "Name":       testServiceName,
        "Repo":       "https://github.com/1tush/docker_test.git",
        "LastCommit": "88782c0d3249d190520335ced2bc41355524708c",
    }
    putServiceConfig(config, t)
    nodes := getNodes(t)
    nodeName := nodes[0].Node
    recorded := test.RunRequest(t, http.StripPrefix("/api", getRestApi().MakeHandler()),
        test.MakeSimpleRequest("POST", fmt.Sprintf("http://localhost/api/services/%s/build?nodes=%s", testServiceName, nodeName), nil))
    data, _ := ioutil.ReadAll(recorded.Recorder.Body)
    t.Logf("%s\n", data)
    recorded.CodeIs(200)
}

func TestSendRestartEvent(t *testing.T) {
    config := map[string]interface{}{
        "Volumes":    []string{},
        "Ports":      []string{"9000:8000"},
        "Env":        []string{"TEST=yes"},
        "Name":       testServiceName,
        "Repo":       "https://github.com/1tush/docker_test.git",
        "LastCommit": "88782c0d3249d190520335ced2bc41355524708c",
    }
    putServiceConfig(config, t)
    nodes := getNodes(t)
    nodeName := nodes[0].Node
    recorded := test.RunRequest(t, http.StripPrefix("/api", getRestApi().MakeHandler()),
        test.MakeSimpleRequest("POST", fmt.Sprintf("http://localhost/api/services/%s/restart?nodes=%s", testServiceName, nodeName), nil))
    data, _ := ioutil.ReadAll(recorded.Recorder.Body)
    t.Logf("%s\n", data)
    recorded.CodeIs(200)
}
