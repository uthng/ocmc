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
)

// GetSwarmServices returns a list of swarm services created in the cluster
func GetSwarmServices (client *client.Client) ([]swarm.Service, error) {
    // Get swarm services
    return client.ServiceList(context.Background(), types.ServiceListOptions{})
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
