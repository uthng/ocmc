package docker

import (
    //"fmt"
    "errors"
    //"strings"
    //"strconv"
    //"time"

    "github.com/rivo/tview"
    "github.com/gdamore/tcell"

    "github.com/uthng/common/utils"

    "github.com/uthng/ocmc/types"
    "github.com/uthng/ocmc/console"
)

/////////////// DECLARATION OF GLOBAL VARIABLES ///////////////////
//var k8sPods             []corev1.PodList
//var swarmTasks          []swarm.Task
//var swarmNodes          []swarm.Node
//var swarmNetworks       []docker_types.NetworkResource
var lastSelectedMenu    string
var orderedMenus        []string
//var nodeClients         []types.NodeClient

// NewModuleDocker initializes a new module for swarm cluster.
//
// It defines functions to setup layout and menu for modules
func NewModuleK8s() *types.Module {
    module := &types.Module {
        Name: "k8s",
        Version: "0.1",
        Description: "Kubernetes 1.9",

        Layout: setupLayoutModule,
        Menus: map[string]types.Menu {
            "pods": types.Menu {
                Name: "pods",
                Layout: setupLayoutPod,
                Close: clearLayoutPod,
                Focus: setFocusPod,
            },
            //"nodes": types.Menu {
                //Name: "nodes",
                //Layout: setupLayoutNode,
                //Close: clearLayoutNode,
                //Focus: setFocusNode,
            //},

        },
    }

    return module
}

// setupLayoutModule sets the layout for the module with a menu list for actions
// and a detail zone
func setupLayoutModule(container string, page *console.Page) error {
    // Setup layout for the 1st elem of menu list
    // Because tview.List does not provide a function to get current selected item
    // so we force to menu "services"
    data, _ := page.Data.(*types.PageClusterData)

    orderedMenus = utils.GetMapSortedKeys(data.Module.Menus, false).([]string)
    if orderedMenus == nil {
        return errors.New("Cannot sort menu keys")
    }

    listMenu := setupListMenu(page)

    err := page.AddItem("cluster", "list_menu", listMenu, 0, 1, false)
    if err != nil {
        return err
    }

    // Add container for cluster
    err = page.AddContainer("cluster", "details", tview.FlexColumn, 0, 10, false)
    if err != nil {
        return err
    }

    // Show details of 1st menu
    data.Module.Menus[orderedMenus[0]].Layout("details", page)
    lastSelectedMenu = orderedMenus[0]

    //details, _ := page.GetContainer("details")
    //details.SetBorder(true)

    //tableService := createTableService(page)
    //err = page.AddItem("details", "table_services", tableService, 0, 1, false)
    //if err != nil {
        //return nil, err
    //}
    return nil
}

// setupListMenu loads a list of available actions and handles
// all key events
func setupListMenu(page *console.Page) *tview.List {
    data, _ := page.Data.(*types.PageClusterData)

    // List menu
    list := tview.NewList().ShowSecondaryText(false)
    list.SetBorder(true).SetTitle("Menu")

    // Only forcing menu display in order that we want
    for _, item := range orderedMenus {
        list.AddItem(item, "", 0, nil)
    }

    // When clicking on an item
    list.SetSelectedFunc(func(i int, menuName string, t string, s rune) {
        if lastSelectedMenu != "" {
            if data.Module.Menus[lastSelectedMenu].Close != nil {
                data.Module.Menus[lastSelectedMenu].Close("details", page)
            }
        }
        data.Module.Menus[menuName].Layout("details", page)
        lastSelectedMenu = menuName
    })

    // Modify certain key events before forwarding others to default handler
    list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        //fmt.Println("Key pressed")
        switch event.Key() {
        case tcell.KeyEsc:
            list, _ := page.GetElemList("list_clusters")
            data.App.SetFocus(list)
            return nil
         case tcell.KeyTab:
            data.Module.Menus[lastSelectedMenu].Focus(page)
            return nil

        }
        return event
    })


    return list
}
