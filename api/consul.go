package api

import (
    porterConsul "github.com/da4nik/porter/consul"
    consulapi "github.com/hashicorp/consul/api"
)

func Services() (services map[string][]string, err error) {
    services, _, err = porterConsul.GetApi().Client().Catalog().Services(nil)
    delete(services, "consul")
    return
}

func Nodes() (nodes []*consulapi.Node, err error) {
    nodes, _, err = porterConsul.GetApi().Client().Catalog().Nodes(nil)
    return
}

func NodeServices(nodeName string) (services map[string]*consulapi.AgentService, err error) {
    node, _, err := porterConsul.GetApi().Client().Catalog().Node(nodeName, nil)
    if err != nil {
        return
    }
    services = node.Services
    if err != nil {
        return
    }
    delete(services, "consul")
    return
}

func fireEvent(serviceName string, nodes []string, creator func(string, string, string) porterConsul.Event) error {
    var (
        event         porterConsul.Event
        serviceFilter string
    )
    _, err := porterConsul.GetServiceConfig(serviceName)
    if err != nil {
        return err
    }
    if len(nodes) == 0 {
        nodes = []string{""}
        serviceFilter = serviceName
    }
    for _, nodeName := range nodes {
        event = creator(serviceName, nodeName, serviceFilter)
        err = event.Fire()
        if err != nil {
            return err
        }
    }
    return nil
}

func FireBuildEvent(serviceName string, nodes []string) error {
    return fireEvent(serviceName, nodes, porterConsul.NewBuildImageEvent)
}

func FireRestartEvent(serviceName string, nodes []string) error {
    return fireEvent(serviceName, nodes, porterConsul.NewRestartContainerEvent)
}
