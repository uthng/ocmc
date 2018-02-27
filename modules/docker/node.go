package docker

import (
    "fmt"
    //"errors"
    "strings"
    "strconv"
    "time"
    //"sort"

    "github.com/rivo/tview"
    "github.com/gdamore/tcell"

    //"golang.org/x/net/context"
    "github.com/docker/docker/client"
    //docker_types "github.com/docker/docker/api/types"
    //"github.com/docker/docker/api/types/swarm"

    "github.com/uthng/ocmc/types"
    "github.com/uthng/ocmc/console"
    page_console "github.com/uthng/ocmc/pages/console"
    "github.com/uthng/ocmc/common/config"
    "github.com/uthng/ocmc/common/docker"
)

// setupLayoutNodes initializes zone containing different elements of
// service
func setupLayoutNode(container string, page *console.Page) error {
    // Set direction for each menu
    page.SetContainerDirection(container, tview.FlexRow)

    err := setupTableNodes(container, page)
    if err != nil {
        return err
    }

    return nil
}

func clearLayoutNode(container string, page *console.Page) {
    //data, _ := page.Data.(*types.PageClusterData)
    err := page.RemoveItem(container, "table_nodes")
    if err != nil {
        fmt.Println(err)
    }

    err = page.RemoveItem(container, "table_containers")
    if err != nil {
        fmt.Println(err)
    }

}

// setFocusNode set focus  manually on the first element of detail zone
func setFocusNode(page *console.Page) {
    data, _ := page.Data.(*types.PageClusterData)

    // Check if table already exists. If not, create it. Otherwise reuse it
    table, err := page.GetElemTable("table_nodes")
    if err == nil {
        data.App.SetFocus(table)
    }
}


// setupTableNodes initializes a table contaning nodes and
// handles key event for navigation
func setupTableNodes(container string, page *console.Page) error {
    var tableNodes *tview.Table

    data, _ := page.Data.(*types.PageClusterData)

    // Check if table already exists. If not, create it. Otherwise reuse it
    tableNodes, err := page.GetElemTable("table_nodes")
    if err != nil {
        // Set column Clusters
        tableNodes = tview.NewTable()
        tableNodes.SetBorders(false)
        tableNodes.SetBorder(true).SetBorderPadding(0, 0, 0, 0).SetTitle("Nodes")
        tableNodes.SetSeparator(tview.GraphicsVertBar)
        tableNodes.SetSelectable(true, false)

        err = page.AddItem(container, "table_nodes", tableNodes, 0, 1, false)
        if err != nil {
            return err
        }
    }

    // Clear and draw table header
    tableNodes.Clear()
    tableNodes.SetCell(0, 0, &tview.TableCell{Text: "ID", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true})
    tableNodes.SetCell(0, 1, &tview.TableCell{Text: "Hostname", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true})
    tableNodes.SetCell(0, 2, &tview.TableCell{Text: "Status", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true})
    tableNodes.SetCell(0, 3, &tview.TableCell{Text: "Availability", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true})
    tableNodes.SetCell(0, 4, &tview.TableCell{Text: "Role", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true})
    tableNodes.SetCell(0, 5, &tview.TableCell{Text: "Manager status", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true})
    tableNodes.SetCell(0, 6, &tview.TableCell{Text: "Address", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true})

    // Get swarm nodes
    swarmNodes, err = docker.GetSwarmNodes(data.Module.Client.(*client.Client))
    if err != nil {
        fmt.Println(err)
        return err
    }

    // Build service table
    for i, node := range swarmNodes {
        tableNodes.SetCell(i+1, 0, &tview.TableCell{Text: node.ID, Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 100 })
        tableNodes.SetCell(i+1, 1, &tview.TableCell{Text: node.Description.Hostname, Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 100 })
        tableNodes.SetCell(i+1, 2, &tview.TableCell{Text: string(node.Status.State), Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 100 })
        tableNodes.SetCell(i+1, 3, &tview.TableCell{Text: string(node.Spec.Availability), Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 100 })
        tableNodes.SetCell(i+1, 4, &tview.TableCell{Text: string(node.Spec.Role), Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 100 })

        managerStatus := ""
        nodeAddress := node.Status.Addr
        if node.ManagerStatus != nil {
            if node.ManagerStatus.Leader {
                managerStatus = "leader"
                nodeAddress = strings.Split(node.ManagerStatus.Addr, ":")[0]
            } else {
                managerStatus = string(node.ManagerStatus.Reachability)
            }
        }
        tableNodes.SetCell(i+1, 5, &tview.TableCell{Text: managerStatus, Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 100 })
        tableNodes.SetCell(i+1, 6, &tview.TableCell{Text: nodeAddress, Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 100 })
    }

    // Point to 1st elem of service table
    setupTableNodeContainers(tableNodes.GetCell(1, 6).Text, container, page)

    // Handle Enter key event on each service
    tableNodes.SetSelectedFunc(func(row int, column int) {
        setupTableNodeContainers(tableNodes.GetCell(row, 6).Text, container, page)
    })

    // Handle other event key than Enter
    tableNodes.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        //fmt.Println("Key pressed")
        switch event.Key() {
        case tcell.KeyEsc:
            list, _ := page.GetElemList("list_menu")
            data.App.SetFocus(list)
            return nil
        case tcell.KeyTab:
            table, _ := page.GetElemTable("table_containers")
            data.App.SetFocus(table)
            return nil
        case tcell.KeyF5:
            // Get current selected row
            row, _ := tableNodes.GetSelection()
            // Get the server ip of the current row
            host := tableNodes.GetCell(row, 6).Text
            // Check if a node client is already initialized.
            // If yes, use it. Otherwise, initialize a new one
            nodeClient, err := getNodeClient(host, data)
            if err != nil {
                return nil
            }
            // New data console
            dataConsole := &types.PageConsoleData{
                PageName: host,
                Node: &nodeClient,
                App: data.App,
            }

            pageConsole, err := page_console.NewPageConsole(dataConsole)
            if err != nil {
                return nil
            }
            data.App.GetPages().AddPage(host, pageConsole, true, true)
            return nil
        }

        return event
    })

    return nil
}

// setupTableServiceContainers create a table for containers' information deployed
// related to the service and handles key events for navigation
func setupTableNodeContainers(server string, container string, page *console.Page) error {
    var table *tview.Table

    data, _ := page.Data.(*types.PageClusterData)

    table, err := page.GetElemTable("table_containers")
    if err != nil {
        // Set column Clusters
        table = tview.NewTable()
        table.SetBorders(false)
        table.SetBorder(true).SetBorderPadding(0, 0, 0, 0).SetTitle("Containers")
        table.SetSeparator(tview.GraphicsVertBar)
        table.SetSelectable(true, false)

        err = page.AddItem(container, "table_containers", table, 0, 1, false)
        if err != nil {
            return err
        }
    }

    // Clear & rebuild attribute table
    table.Clear()
    table.SetCell(0, 0, &tview.TableCell{Text: "ID", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true})
    table.SetCell(0, 1, &tview.TableCell{Text: "Name", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true})
    table.SetCell(0, 2, &tview.TableCell{Text: "State", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true})
    table.SetCell(0, 3, &tview.TableCell{Text: "Status", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true})
    table.SetCell(0, 4, &tview.TableCell{Text: "Image", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true})
    table.SetCell(0, 5, &tview.TableCell{Text: "Created", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true})
    table.SetCell(0, 6, &tview.TableCell{Text: "Ports", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true})
    table.SetCell(0, 7, &tview.TableCell{Text: "Command", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true})

    // Check if a node client is already initialized.
    // If yes, use it. Otherwise, initialize a new one
    nodeClient, err := getNodeClient(server, data)
    if err != nil {
        return err
    }

    nodeClients = append(nodeClients, nodeClient)

    // Get containers
    containers, err := docker.GetContainers(nodeClient.Client.(*client.Client))
    if err != nil {
        return err
    }

    for i, c := range containers {
        table.SetCell(i+1, 0, &tview.TableCell{Text: c.ID, Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 30 })
        table.SetCell(i+1, 1, &tview.TableCell{Text: strings.Trim(c.Names[0], "/"), Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 30 })
        table.SetCell(i+1, 2, &tview.TableCell{Text: c.State, Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 100 })
        table.SetCell(i+1, 3, &tview.TableCell{Text: c.Status, Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 100 })
        table.SetCell(i+1, 4, &tview.TableCell{Text: strings.Split(c.Image, "@")[0], Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 100 })
        table.SetCell(i+1, 5, &tview.TableCell{Text: time.Unix(c.Created, 0).UTC().Format(time.UnixDate), Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 100 })

        ports := ""
        for _, port := range c.Ports {
            if port.PublicPort > 0  {
                ports = ports + " " + strconv.Itoa(int(port.PublicPort)) + "->" + strconv.Itoa(int(port.PrivatePort)) + "/" + port.Type
            } else {
                ports = ports + " " + strconv.Itoa(int(port.PrivatePort)) + "/" + port.Type
            }
        }

        table.SetCell(i+1, 6, &tview.TableCell{Text: ports, Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 100 })
        table.SetCell(i+1, 7, &tview.TableCell{Text: c.Command, Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 100 })
    }

    // Handle other event key than Enter
    table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        //fmt.Println("Key pressed")
        switch event.Key() {
        case tcell.KeyEsc:
            t, _ := page.GetElemTable("table_nodes")
            data.App.SetFocus(t)
            return nil
         //case tcell.KeyTab:
            //table, _ := page.GetElemTable("table_attributes")
            //data.App.SetFocus(table)
            //return nil
        }
        return event
    })

    return nil
}

// getNodeClient searchs to see if there is a node client corresponding to
// the given server in the global variable nodeClients of this module
// If yes, use it. Otherwise, create a new one
func getNodeClient(server string, data *types.PageClusterData) (types.NodeClient, error) {
    client := types.NodeClient{}

    for _, c := range nodeClients {
        if c.Config.Host == server {
            client = c
        }
    }

    if client.Client == nil {
        // Prepare config connection
        connConfig := config.GetClusterConfig(data.PageName, data).Config
        connConfig.Host = server

        // Initialize new node client
        c, err := docker.NewDockerClient(connConfig)
        if err != nil {
            return client, err
        }

        // Init new client
        client.Client = c
        client.Config = connConfig
    }

    return client, nil
}
