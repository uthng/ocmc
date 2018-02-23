package types

import (

    "github.com/uthng/ocmc/console"
)


type ClusterConfig struct {
    Name                string
    Config              ConnConfig
}


type PageClusterData struct {
    PageName            string
    Configs             []ClusterConfig

    //Client              interface{}
    Module              *Module

    App                 *console.App
}
