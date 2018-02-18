package docker

import (
    //"fmt"
    "errors"

    "golang.org/x/net/context"
    "github.com/moby/moby/client"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/swarm"

)
var (
    ErrNetworkIDNotFound = errors.New("Network ID not found")
    ErrServiceIDNotFound = errors.New("Service ID not found")
    ErrServiceNameNotFound = errors.New("Service Name not found")
    ErrNodeIDNotFound = errors.New("Node ID not found")
)

// GetSwarmServices returns a list of swarm services created in the cluster
func GetSwarmServices (client *client.Client) ([]swarm.Service, error) {
    // Get swarm services
    return client.ServiceList(context.Background(), types.ServiceListOptions{})
}

// FindServiceByID searchs and returns the swarm service corresponding to
// the given ID
func FindSwarmServiceByID (id string, services []swarm.Service) (*swarm.Service, error) {
    for _, service := range services {
        if service.ID == id {
            return &service, nil
        }
    }

    return nil, ErrServiceIDNotFound
}

// FindServiceByName searchs and returns the swarm service corresponding to
// the given name
func FindSwarmServiceByName (name string, services []swarm.Service) (*swarm.Service, error) {
    for _, service := range services {
        if service.Spec.Annotations.Name == name {
            return &service, nil
        }
    }

    return nil, ErrServiceNameNotFound
}

// GetNetworks return a list of network defined in the cluster
func GetNetworks (client *client.Client) ([]types.NetworkResource, error) {
    return client.NetworkList(context.Background(), types.NetworkListOptions{})
}

// FindNetworkByID searchs and returns the NetworkResource corresponding to
// the given ID
func FindNetworkByID (id string, networks []types.NetworkResource) (*types.NetworkResource, error) {
    for _, net := range networks {
        if net.ID == id {
            return &net, nil
        }
    }

    return nil, ErrNetworkIDNotFound
}

// GetSwarmTasks return a list of swarm containers in the cluster
func GetSwarmTasks (client *client.Client) ([]swarm.Task, error) {
    return client.TaskList(context.Background(), types.TaskListOptions{})
}

// FindSwarmTasksByServiceID searchs and returns the tasks
// related to the given service ID
func FindSwarmTasksByServiceID (serviceId string, tasks []swarm.Task) ([]swarm.Task) {
    var result []swarm.Task

    for _, task := range tasks {
        if task.ServiceID == serviceId {
            result = append(result, task)
        }
    }

    return result
}

// GetSwarmNodes return a list of swarm nodes in the cluster
func GetSwarmNodes (client *client.Client) ([]swarm.Node, error) {
    return client.NodeList(context.Background(), types.NodeListOptions{})
}

// FindNodeByServiceID returns the node corresponding to the given ID
func FindSwarmNodeByID (id string, nodes []swarm.Node) (*swarm.Node, error) {
    for _, node := range nodes {
        if node.ID == id {
            return &node, nil
        }
    }

    return nil, ErrNodeIDNotFound
}

