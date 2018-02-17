package docker

import (
    //"fmt"
    //"errors"
    "strings"
    "strconv"
    "time"

    "github.com/rivo/tview"
    "github.com/gdamore/tcell"

    //"golang.org/x/net/context"
    docker "github.com/moby/moby/client"
    //docker_types "github.com/docker/docker/api/types"
    //"github.com/docker/docker/api/types/swarm"

    "github.com/uthng/ocmc/types"
    "github.com/uthng/ocmc/console"
    //"github.com/uthng/ocmc/pages"
)

func setupLayoutService(container string, page *console.Page) error {
    err := setupTableService(container, page)
    if err != nil {
        return err
    }

    return nil
}


func setupTableService(container string, page *console.Page) error {
    var tableService *tview.Table

    data, _ := page.GetData().(*types.PageClusterData)

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
    }

    // Clear and draw table header
    tableService.Clear()
    tableService.SetCell(0, 0, &tview.TableCell{Text: "Name", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true})

    // Get swarm services
    swarmServices, err = GetSwarmServices(data.Module.Client.(*docker.Client))
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
    setupTableServiceAttribute(tableService.GetCell(1, 0).Text, container, page)

    // Handle Enter key event on each service
    tableService.SetSelectedFunc(func(row int, column int) {
        setupTableServiceAttribute(tableService.GetCell(row, column).Text, container, page)
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

        }
        return event
    })

    return nil
}

func setupTableServiceAttribute(service string, container string, page *console.Page) error {
    var table *tview.Table
    var attributes = make(map[string]string)
    var keyAttributes []string

    data, _ := page.GetData().(*types.PageClusterData)

    table, err := page.GetElemTable("table_attributes")
    if err != nil {
        // Set column Clusters
        table = tview.NewTable()
        table.SetBorders(false)
        table.SetBorder(true).SetBorderPadding(0, 0, 0, 0).SetTitle("Attributes")
        table.SetSeparator(tview.GraphicsVertBar)
        table.SetSelectable(true, false)

        err = page.AddItem(container, "table_attributes", table, 0, 3, false)
        if err != nil {
            return err
        }
    }

    // Clear & rebuild attribute table
    table.Clear()
    table.SetCell(0, 0, &tview.TableCell{Text: "Attribute", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true}).
    SetCell(0, 1, &tview.TableCell{Text: "Value", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true})

    // Get swarm networks
    swarmNetworks, err = GetNetworks(data.Module.Client.(*docker.Client))
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
                network, err := FindNetworkByID(vip.NetworkID, swarmNetworks)
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
         //case tcell.KeyTab:
            //table, _ := page.GetElemTable("table_attributes")
            //data.App.SetFocus(table)
            //return nil

        }
        return event
    })

    return nil
}
