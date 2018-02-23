package types

import (

    "github.com/uthng/ocmc/console"
)


type PageConsoleData struct {
    PageName            string

    Node                *NodeClient
    App                 *console.App
}

