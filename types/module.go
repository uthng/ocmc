package types

import (
    //    "fmt"
    //"errors"

    "github.com/uthng/ocmc/console"

)

/////// DECLARATION OF ALL TYPES /////////////
//type CmdResult struct {
    //Error       error
    //Result      map[string]interface{}
//}

//type CmdFunc func(map[string]interface{}) *CmdResult

type LayoutFunc func(container  string, page *console.Page) error
type CloseFunc func(container  string, page *console.Page)
type FocusFunc func(page *console.Page)

type Menu struct {
    Name        string

    //Cmd         CmdFunc
    Layout      LayoutFunc
    Close       CloseFunc
    Focus       FocusFunc
}

// Struct for a command of a given module
type Module struct {
    Name        string
    Version     string
    Description string

    Client      interface{}
    Layout      LayoutFunc
    Menus       map[string]Menu
}

// Nodeclient contains name and client connection to a node
// in the cluster
// This struct can be used to execute a command on a container
// on the node
type NodeClient struct {
    // Server host
    Host        string
    // type: docker or kubernetes
    Type        string
    // CLient: docker or kubernetes or other stuffs
    Client      interface{}
}
