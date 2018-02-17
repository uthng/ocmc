package types

import (

    "github.com/uthng/ocmc/console"
)


type ClusterConfig struct {
    Name        string
    Host        string
    Type        string
    AuthType    string
    Auth        ClusterAuth
}

type ClusterAuth struct {
    Type        string
    Ca          string
    Client      string
    ClientKey   string

    SshKey      string

    Username    string
    Passord     string
}

type PageClusterData struct {
    PageName            string
    Configs             []ClusterConfig

    //Client              interface{}
    Module              *Module

    App                 *console.App
}

