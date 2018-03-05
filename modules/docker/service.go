package docker

import (
    "fmt"
    //"errors"
    "strings"
    "strconv"
    "time"
    "sort"

    "github.com/rivo/tview"
    "github.com/gdamore/tcell"

    "github.com/uthng/common/docker"

    "github.com/uthng/ocmc/types"
    "github.com/uthng/ocmc/console"
)

// setupLayoutService initializes zone containing different elements of
// service
func setupLayoutService(container string, page *console.Page) error {
    // Set default direction for the current menu
    page.SetContainerDirection(container, tview.FlexColumn)

    err := setupTableService(container, page)
    if err != nil {
        return err
    }

    return nil
}

func clearLayoutService(container string, page *console.Page) {
    //data, _ := page.Data.(*yytypes.PageClusterData)
    err := page.RemoveItem("display_details", "table_attributes")
    if err != nil {
        fmt.Println(err)
    }

    err = page.RemoveItem("display_details", "table_tasks")
    if err != nil {
        fmt.Println(err)
    }

    err = page.RemoveContainer("details", "display_details")
    if err != nil {
        fmt.Println(err)
    }

    err = page.RemoveItem("details", "table_services")
    if err != nil {
        fmt.Println(err)
    }

}

// setFocusNode set focus  manually on the first element of detail zone
func setFocusService(page *console.Page) {
    data, _ := page.Data.(*types.PageClusterData)

    // Check if table already exists. If not, create it. Otherwise reuse it
    tableService, err := page.GetElemTable("table_services")
    if err == nil {
        data.App.SetFocus(tableService)
    }
}

// setupTableService initializes a table contaning services and
// handles key event for navigation
func setupTableService(container string, page *console.Page) error {
    var tableService *tview.Table

    data, _ := page.Data.(*types.PageClusterData)

    // Check if table already exists. If not, create it. Otherwise reuse it
    tableService, err := page.GetElemTable("table_services")
    if err != nil {
        // Set column Clusters
        tableService = tview.NewTable()
        tableService.SetBorders(false)
        tableService.SetBorder(true).SetBorderPadding(0, 0, 0, 0).SetTitle("Services")
        //tableService.SetSeparator(tview.GraphicsVertBar)
        tableService.SetSelectable(true, false)

        err = page.AddItem(container, "table_services", tableService, 0, 1, false)
        if err != nil {
            return err
        }

        // Add new container to include 2 tables: service atttributes and containers
        err = page.AddContainer(container, "display_details", tview.FlexRow, 0, 3, false)
        if err != nil {
            return err
        }

    }

    // Clear and draw table header
    tableService.Clear()
    tableService.SetCell(0, 0, &tview.TableCell{Text: "Name", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true})

    // Get swarm services
    client := data.Module.Client.(*docker.Client)
    swarmServices, err = client.GetSwarmServices(ctx, nil)
    if err != nil {
        return err
    }

    // Build service table
    for i, service := range swarmServices {
        //image := strings.Split(service.Spec.TaskTemplate.ContainerSpec.Image, "@")[0]
        tableService.SetCell(i+1, 0, &tview.TableCell{Text: service.Spec.Annotations.Name, Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 40 })
        //fmt.Println(service.Spec.Annotations.Name)
    }

    // Point to 1st elem of service table
    //tableService.Select(1, 0)
    setupTableServiceAttributes(tableService.GetCell(1, 0).Text, "display_details", page)
    setupTableServiceContainers(tableService.GetCell(1, 0).Text, "display_details", page)

    // Handle Enter key event on each service
    tableService.SetSelectedFunc(func(row int, column int) {
        setupTableServiceAttributes(tableService.GetCell(row, column).Text, "display_details", page)
        setupTableServiceContainers(tableService.GetCell(row, column).Text, "display_details", page)
    })

    // Handle other event key than Enter
    tableService.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        //fmt.Println("Key pressed")
        switch event.Key() {
        case tcell.KeyEsc:
            list, _ := page.GetElemList("list_menu")
            data.App.SetFocus(list)
            return nil
        case tcell.KeyTab:
            table, _ := page.GetElemTable("table_attributes")
            data.App.SetFocus(table)
            return nil
        case tcell.KeyF5:
            setupTableService(container, page)
            return nil
        }
        return event
    })

    return nil
}

// setupTableServiceAttributes create a table containing service attributes
// and handles key events for navigaition
func setupTableServiceAttributes(service string, container string, page *console.Page) error {
    var table *tview.Table
    var attributes = make(map[string]string)
    var keyAttributes []string

    data, _ := page.Data.(*types.PageClusterData)

    table, err := page.GetElemTable("table_attributes")
    if err != nil {
        // Set column Clusters
        table = tview.NewTable()
        table.SetBorders(false)
        table.SetBorder(true).SetBorderPadding(0, 0, 0, 0).SetTitle("Attributes")
        table.SetSeparator(tview.GraphicsVertBar)
        table.SetSelectable(true, false)

        err = page.AddItem(container, "table_attributes", table, 0, 1, false)
        if err != nil {
            return err
        }
    }

    // Clear & rebuild attribute table
    table.Clear()
    table.SetCell(0, 0, &tview.TableCell{Text: "Attribute", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true}).
    SetCell(0, 1, &tview.TableCell{Text: "Value", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true})

    // Get swarm networks
    client := data.Module.Client.(*docker.Client)
    swarmNetworks, err = client.GetNetworks(ctx, nil)
    if err != nil {
        return err
    }

    for _, srv := range swarmServices {
        //fmt.Printf("Service attributes %v\n", serv)
        if srv.Spec.Annotations.Name == service {
            // Build map[string][string]
            attributes["ID"] = srv.ID
            attributes["Image"] = strings.Split(srv.Spec.TaskTemplate.ContainerSpec.Image, "@")[0]
            // Get service mode
            if srv.Spec.Mode.Global == nil {
                attributes["Mode"] = "Replicated"
                attributes["Replicas"] = strconv.FormatUint(*srv.Spec.Mode.Replicated.Replicas, 10)
            } else {
                attributes["Mode"] = "Global"
            }

            // Append key in order
            keyAttributes = append(keyAttributes, "ID", "Image", "Mode", "Replicas")

            // Get Ports
            for i, port := range srv.Endpoint.Ports {
                attributes["Port" + strconv.Itoa(i)] = strconv.Itoa(int(port.PublishedPort)) + "->" + strconv.Itoa(int(port.TargetPort)) + "/" + string(port.Protocol)
                keyAttributes = append(keyAttributes, "Port" + strconv.Itoa(i))
            }

            attributes["CreatedAt"] = srv.Meta.CreatedAt.UTC().Format(time.UnixDate)
            attributes["UpdatedAt"] = srv.Meta.UpdatedAt.UTC().Format(time.UnixDate)

            // Append key in order
            keyAttributes = append(keyAttributes, "CreatedAt", "UpdatedAt")

            // Get virtual ip with network
            for i, vip := range srv.Endpoint.VirtualIPs {
                network, err := client.FindNetworkByID(vip.NetworkID, swarmNetworks)
                if err != nil {
                    attributes["VIP" + strconv.Itoa(i) + " - Network"] = err.Error()
                } else {
                    attributes["VIP" + strconv.Itoa(i) + " - Network"] = network.Name
                }

                attributes["VIP" + strconv.Itoa(i) + " - IP"] = vip.Addr

                // Append key in order
                keyAttributes = append(keyAttributes, "VIP" + strconv.Itoa(i) + " - Network", "VIP" + strconv.Itoa(i) + " - IP")
            }
            //attributes[""] = service.ID
        }
    }

    i := 0
    maxWidthValue := 100
    for _, key := range keyAttributes {
        //image := strings.Split(service.Spec.TaskTemplate.ContainerSpec.Image, "@")[0]
        value, ok := attributes[key]
        if ok == true {
            table.SetCell(i+1, 0, &tview.TableCell{Text: key, Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 30 })
            if len(value) > maxWidthValue {
                table.SetCell(i+1, 1, &tview.TableCell{Text: value[:maxWidthValue], Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: maxWidthValue })
                i = i + 1
                table.SetCell(i+1, 1, &tview.TableCell{Text: value[maxWidthValue:], Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: maxWidthValue })
            } else {
                table.SetCell(i+1, 1, &tview.TableCell{Text: value, Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: maxWidthValue })
            }
            i = i + 1
        }

    }


    // Handle other event key than Enter
    table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        //fmt.Println("Key pressed")
        switch event.Key() {
        case tcell.KeyEsc:
            t, _ := page.GetElemTable("table_services")
            data.App.SetFocus(t)
            return nil
        case tcell.KeyTab:
            table, _ := page.GetElemTable("table_tasks")
            data.App.SetFocus(table)
            return nil
        case tcell.KeyF5:
            setupTableServiceAttributes(service, container, page)
            return nil

        }
        return event
    })

    return nil
}

// setupTableServiceContainers create a table for containers' information deployed
// related to the service and handles key events for navigation
func setupTableServiceContainers(service string, container string, page *console.Page) error {
    var table *tview.Table

    data, _ := page.Data.(*types.PageClusterData)

    table, err := page.GetElemTable("table_tasks")
    if err != nil {
        // Set column Clusters
        table = tview.NewTable()
        table.SetBorders(false)
        table.SetBorder(true).SetBorderPadding(0, 0, 0, 0).SetTitle("Tasks")
        table.SetSeparator(tview.GraphicsVertBar)
        table.SetSelectable(true, false)

        err = page.AddItem(container, "table_tasks", table, 0, 1, false)
        if err != nil {
            return err
        }
    }

    // Clear & rebuild attribute table
    table.Clear()
    table.SetCell(0, 0, &tview.TableCell{Text: "ID", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true}).
          //SetCell(0, 1, &tview.TableCell{Text: "Name", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true}).
          SetCell(0, 1, &tview.TableCell{Text: "Node", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true}).
          SetCell(0, 2, &tview.TableCell{Text: "Desired", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true}).
          SetCell(0, 3, &tview.TableCell{Text: "Current", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true}).
          SetCell(0, 4, &tview.TableCell{Text: "Created", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true}).
          SetCell(0, 5, &tview.TableCell{Text: "Updated", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true})

    // Get swarm networks
    client := data.Module.Client.(*docker.Client)
    swarmTasks, err = client.GetSwarmTasks(ctx, nil)
    if err != nil {
        return err
    }

    // Find sevrice id corresponding to name
    srv, err := client.FindSwarmServiceByName(service, swarmServices)
    if err != nil {
        return err
    }

    // List all tasks / containers related to service
    serviceTasks := client.FindSwarmTasksByServiceID(srv.ID, swarmTasks)

    // Get list of nodes
    swarmNodes, err := client.GetSwarmNodes(ctx, nil)
    if err != nil {
        return err
    }

    // Sort the slice following creation date
    sort.Slice(serviceTasks, func(i, j int) bool { return serviceTasks[i].Meta.CreatedAt.After(serviceTasks[j].Meta.CreatedAt)})

    for i, task := range serviceTasks {
        table.SetCell(i+1, 0, &tview.TableCell{Text: task.ID, Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 30 })
        //table.SetCell(i+1, 1, &tview.TableCell{Text: task.Annotations.Name, Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 30 })

        node, err := client.FindSwarmNodeByID(task.NodeID, swarmNodes)
        if err != nil {
            table.SetCell(i+1, 1, &tview.TableCell{Text: "", Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 30 })
        } else {
            table.SetCell(i+1, 1, &tview.TableCell{Text: node.Description.Hostname, Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 30 })
        }
        table.SetCell(i+1, 2, &tview.TableCell{Text: string(task.DesiredState), Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 30 })
        table.SetCell(i+1, 3, &tview.TableCell{Text: string(task.Status.State), Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 30 })
        table.SetCell(i+1, 4, &tview.TableCell{Text: task.Meta.CreatedAt.UTC().Format(time.UnixDate), Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 30 })
        table.SetCell(i+1, 5, &tview.TableCell{Text: task.Meta.UpdatedAt.UTC().Format(time.UnixDate), Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 30 })
    }

    // Handle other event key than Enter
    table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        //fmt.Println("Key pressed")
        switch event.Key() {
        case tcell.KeyEsc:
            t, _ := page.GetElemTable("table_attributes")
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
