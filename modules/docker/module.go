package docker

import (
    //"fmt"
    "errors"
    //"strings"
    //"strconv"
    //"time"

    "github.com/rivo/tview"
    "github.com/gdamore/tcell"

    "golang.org/x/net/context"
    //docker "github.com/moby/moby/client"
    docker_types "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/swarm"

    "github.com/uthng/common/utils"
    "github.com/uthng/common/docker"

    "github.com/uthng/ocmc/types"
    "github.com/uthng/ocmc/console"
)

/////////////// DECLARATION OF GLOBAL VARIABLES ///////////////////
var swarmServices       []swarm.Service
var swarmTasks          []swarm.Task
var swarmNodes          []swarm.Node
var swarmNetworks       []docker_types.NetworkResource
var lastSelectedMenu    string
var orderedMenus        []string
var nodeClients         []types.NodeClient

var ctx = context.Background()

// NewModuleDocker initializes a new module for swarm cluster.
//
// It defines functions to setup layout and menu for modules
func NewModuleDocker(config types.ConnConfig) (*types.Module, error) {
    module := &types.Module {
        Name: "docker",
        Version: "0.1",
        Description: "Docker and Swarm",

        Layout: setupLayoutModule,
        Menus: map[string]types.Menu {
            "services": types.Menu {
                Name: "services",
                Layout: setupLayoutService,
                Close: clearLayoutService,
                Focus: setFocusService,
            },
            "nodes": types.Menu {
                Name: "nodes",
                Layout: setupLayoutNode,
                Close: clearLayoutNode,
                Focus: setFocusNode,
            },

        },
    }

    module.Client, err = newDockerClient(config)
    if err != nil {
        return nil, errors.New("Cannot initialize docker client")
    }

    return module, nil
}

// NewDockerClient initializes a docker client to remote cluster
// following authentication configuration
func NewDockerClient(config ocmc_types.ConnConfig) (interface{}, error) {
    var client interface{}

    if config.Auth.Type == "ssh" {
        if config.Auth.Kind == "key" {
            sshConfig, err := ssh.NewClientConfigWithKeyFile(config.Auth.Username, config.Auth.SshKey, "", 0, false)
            if err != nil {
                return nil, err
            }

            client, err = docker.NewSSHClient(config.Host + ":" + strconv.Itoa(config.Port), "/var/run/docker.sock", "1.35", sshConfig.ClientConfig)
            if err != nil {
                return nil, err
            }
        }
    }
    return client, nil
}

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
        //if menuName == "services" {
            //client := data.Client.(*client.Client)
            //services, err := client.ServiceList(context.Background(), types.ServiceListOptions{})
            //if err != nil {
                //fmt.Println(err)
            //}

            //table, err := page.GetElemTable("table_services")
            //for i, service := range services {
                ////fmt.Println(i)
                //image := strings.Split(service.Spec.TaskTemplate.ContainerSpec.Image, "@")[0]
                //table.SetCell(i+1, 0, &tview.TableCell{Text: service.ID, Align: tview.AlignCenter, Color: tcell.ColorYellow, MaxWidth: 20 }).
                      //SetCell(i+1, 1, &tview.TableCell{Text: service.Spec.Annotations.Name, Align: tview.AlignCenter, Color: tcell.ColorYellow}).

                      ////SetCell(i+1, 2, &tview.TableCell{Text: service.Spec.Mode.Replicated.Replicas, Align: tview.AlignCenter, Color: tcell.ColorYellow}).
                      //SetCell(i+1, 2, &tview.TableCell{Text: "replicas", Align: tview.AlignCenter, Color: tcell.ColorYellow}).
                      //SetCell(i+1, 3, &tview.TableCell{Text: image, Align: tview.AlignLeft, Color: tcell.ColorYellow, MaxWidth: 50})
                ////fmt.Println(service.Spec.Annotations.Name)
            //}
            //table.Clear()
        //}
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
