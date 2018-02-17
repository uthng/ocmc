package types

import (
    //    "fmt"
    //"errors"

    "github.com/uthng/ocmc/console"

)

/////// DECLARATION OF ALL TYPES /////////////
type CmdResult struct {
    Error       error
    Result      map[string]interface{}
}

type CmdFunc func(map[string]interface{}) *CmdResult

type LayoutFunc func(container  string, page *console.Page) error

type Menu struct {
    Name        string

    Cmd         CmdFunc
    Layout      LayoutFunc
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
