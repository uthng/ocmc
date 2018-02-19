package docker

import (
    "fmt"
    //"errors"
    "strings"
    //"strconv"
    //"time"
    //"sort"

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
    //data, _ := page.GetData().(*types.PageClusterData)
    err := page.RemoveItem(container, "table_nodes")
    if err != nil {
        fmt.Println(err)
    }

    //err = page.RemoveItem("details", "table_services")
    //if err != nil {
        //fmt.Println(err)
    //}

}

// setFocusNode set focus  manually on the first element of detail zone
func setFocusNode(page *console.Page) {
    data, _ := page.GetData().(*types.PageClusterData)

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

    data, _ := page.GetData().(*types.PageClusterData)

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
    swarmNodes, err = GetSwarmNodes(data.Module.Client.(*docker.Client))
    if err != nil {
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
    //tableNodes.Select(1, 0)
    //setupTableNodeTasks(tableNodes.GetCell(1, 0).Text, container, page)
    //setupTableNodesContainers(tableNodes.GetCell(1, 0).Text, "display_details", page)

    // Handle Enter key event on each service
    tableNodes.SetSelectedFunc(func(row int, column int) {
        //setupTableNodeTasks(tableNodes.GetCell(row, column).Text, container, page)
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
            table, _ := page.GetElemTable("table_nodes")
            data.App.SetFocus(table)
            return nil

        }
        return event
    })

    return nil
}

// setupTableServiceContainers create a table for containers' information deployed
// related to the service and handles key events for navigation
//func setupTableNodeContainers(service string, container string, page *console.Page) error {
    //var table *tview.Table

    //data, _ := page.GetData().(*types.PageClusterData)

    //table, err := page.GetElemTable("table_containers")
    //if err != nil {
        //// Set column Clusters
        //table = tview.NewTable()
        //table.SetBorders(false)
        //table.SetBorder(true).SetBorderPadding(0, 0, 0, 0).SetTitle("Containers")
        //table.SetSeparator(tview.GraphicsVertBar)
        //table.SetSelectable(true, false)

        //err = page.AddItem(container, "table_containers", table, 0, 1, false)
        //if err != nil {
            //return err
        //}
    //}

    //// Clear & rebuild attribute table
    //table.Clear()
    //table.SetCell(0, 0, &tview.TableCell{Text: "ID", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true}).
          ////SetCell(0, 1, &tview.TableCell{Text: "Name", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true}).
          //SetCell(0, 1, &tview.TableCell{Text: "Node", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true}).
          //SetCell(0, 2, &tview.TableCell{Text: "Desired", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true}).
          //SetCell(0, 3, &tview.TableCell{Text: "Current", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true}).
          //SetCell(0, 4, &tview.TableCell{Text: "Created", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true}).
          //SetCell(0, 5, &tview.TableCell{Text: "Updated", Align: tview.AlignCenter, Color: tcell.ColorYellow, NotSelectable: true})

    //// Get swarm networks
    //swarmTasks, err = GetSwarmTasks(data.Module.Client.(*docker.Client))
    //if err != nil {
        //return err
    //}

    //// Find sevrice id corresponding to name
    //srv, err := FindSwarmServiceByName(service, swarmServices)
    //if err != nil {
        //return err
    //}

    //// List all tasks / containers related to service
    //serviceTasks := FindSwarmTasksByServiceID(srv.ID, swarmTasks)

    //// Get list of nodes
    //swarmNodes, err := GetSwarmNodes(data.Module.Client.(*docker.Client))
    //if err != nil {
        //return err
    //}

    //// Sort the slice following creation date
    //sort.Slice(serviceTasks, func(i, j int) bool { return serviceTasks[i].Meta.CreatedAt.After(serviceTasks[j].Meta.CreatedAt)})

    //for i, task := range serviceTasks {
        //table.SetCell(i+1, 0, &tview.TableCell{Text: task.ID, Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 30 })
        ////table.SetCell(i+1, 1, &tview.TableCell{Text: task.Annotations.Name, Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 30 })

        //node, err := FindSwarmNodeByID(task.NodeID, swarmNodes)
        //if err != nil {
            //table.SetCell(i+1, 1, &tview.TableCell{Text: "", Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 30 })
        //} else {
            //table.SetCell(i+1, 1, &tview.TableCell{Text: node.Description.Hostname, Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 30 })
        //}
        //table.SetCell(i+1, 2, &tview.TableCell{Text: string(task.DesiredState), Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 30 })
        //table.SetCell(i+1, 3, &tview.TableCell{Text: string(task.Status.State), Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 30 })
        //table.SetCell(i+1, 4, &tview.TableCell{Text: task.Meta.CreatedAt.UTC().Format(time.UnixDate), Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 30 })
        //table.SetCell(i+1, 5, &tview.TableCell{Text: task.Meta.UpdatedAt.UTC().Format(time.UnixDate), Align: tview.AlignLeft, Color: tcell.ColorWhite, MaxWidth: 30 })
    //}

    //// Handle other event key than Enter
    //table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        ////fmt.Println("Key pressed")
        //switch event.Key() {
        //case tcell.KeyEsc:
            //t, _ := page.GetElemTable("table_attributes")
            //data.App.SetFocus(t)
            //return nil
         ////case tcell.KeyTab:
            ////table, _ := page.GetElemTable("table_attributes")
            ////data.App.SetFocus(table)
            ////return nil

        //}
        //return event
    //})

    //return nil
//}
