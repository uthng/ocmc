package types

import (

)

// ConnConfig contains informations for connection to host
type ConnConfig struct {
    // Host: ip or domain name
    Host        string
    // Port
    Port        int
    // Type: docker (swarm) or k8s (kubernetes)
    Type        string
    // Auth: connection authentication
    Auth        ConnAuth
}

// ConnAuth contains informations about authentication type
type ConnAuth struct {
    // Type: ssh, cert
    Type        string
    // Kind: key, vault, password
    Kind        string
    // Certficates: file path or from vault
    Ca          string
    Client      string
    ClientKey   string

    // SSH key: file path or from vault
    SshKey      string

    // Username & password when Kind is password
    Username    string
    Passord     string
}

// Nodeclient use config of type ClusterConfig to connect to the server.
// in the cluster
// This struct can be used to execute a command on a container
// on the node
type NodeClient struct {
    // Config: configuration of connection
    Config      ConnConfig
    // CLient: docker or kubernetes or other stuffs
    Client      interface{}
}
