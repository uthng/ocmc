package pages

import (
    "fmt"
    //"strings"

    "github.com/rivo/tview"
    "github.com/gdamore/tcell"

    //"golang.org/x/net/context"

    // docker
    //"github.com/moby/moby/client"
    //"github.com/docker/docker/api/types"
    //"github.com/docker/docker/api/types/filters"
    //"github.com/docker/docker/api/types/swarm"

    "github.com/uthng/common/ssh"
    "github.com/uthng/common/docker"

    "github.com/uthng/ocmc/console"
    "github.com/uthng/ocmc/types"
    module_docker "github.com/uthng/ocmc/modules/docker"
)


func NewPageCluster(data *types.PageClusterData) (*console.Page, error) {
    page := console.NewPage()

    // Set page data
    page.SetData(data)

    // Setup Title
    title := tview.NewTextView()
    fmt.Fprintf(title, "%s", data.PageName)
    title.SetTextAlign(tview.AlignCenter)
    title.SetBorder(true)

    err := page.SetContainerDirection("root", tview.FlexRow)
    if err != nil {
        return nil, err
    }

    // Add title
    err = page.AddItem("root", "title", title, 0, 1, false)
    if err != nil {
        return nil, err
    }

    // Add container for zone project
    err = page.AddContainer("root", "main", tview.FlexColumn, 0, 10, true)
    if err != nil {
        return nil, err
    }

    // Set column Clusters
    listClusters := createListCluster(page)
    //listClusters := tview.NewList().ShowSecondaryText(false)
    //listClusters.SetBorder(true).SetTitle("Clusters")
    //clusters.SetBackgroundColor(tcell.ColorDarkViolet)

    err = page.AddItem("main", "list_clusters", listClusters, 0, 1, true)
    if err != nil {
        return nil, err
    }

    // Add container for cluster
    err = page.AddContainer("main", "cluster", tview.FlexColumn, 0, 7, false)
    if err != nil {
        return nil, err
    }

    clusterConfig := getClusterConfig(data.PageName, data)
    if clusterConfig.Type == "docker" {
        data.Module = module_docker.NewModuleDocker()
        data.Module.Client, err = initSwarmClient(clusterConfig)
        if err != nil {
            fmt.Println(err)
            return nil, err
        }

        data.Module.Layout("cluster", page)
        //fmt.Println("data %v\n", data.Module)
        //fmt.Println("data %v\n", page.GetData().(*types.PageClusterData).Module)

    }

    // Setup menu for project
    //listMenu := createListMenu(page)

    //err = page.AddItem("cluster", "list_menu", listMenu, 0, 1, false)
    //if err != nil {
        //return nil, err
    //}

    // Setup zone showing project details
    //details := tview.NewBox()
    //details.SetBorder(true).SetTitle("Details")

    //err = page.AddItem("cluster", "details", details, 0, 5, false)
    //if err != nil {
        //return nil, err
    //}

    // Add container for cluster
    //err = page.AddContainer("cluster", "details", tview.FlexColumn, 0, 10, false)
    //if err != nil {
        //return nil, err
    //}

    //tableService := createTableService(page)
    //err = page.AddItem("details", "table_services", tableService, 0, 1, false)
    //if err != nil {
        //return nil, err
    //}


    // Init cluster client
    //err = initClusterClient(page)
    //if err != nil {
        //return nil, err
    //}

    return page, nil
}


///////////// PRIVATE FUNCTIONS ///////////////////
func createListCluster(page *console.Page) *tview.List {
    data := page.GetData().(*types.PageClusterData)

    // Set column Clusters
    list := tview.NewList().ShowSecondaryText(false)
    list.SetBorder(true).SetTitle("Clusters")
    //clusters.SetBackgroundColor(tcell.ColorDarkViolet)

    // Populate list content
    for _, conf := range data.Configs {
        list.AddItem(conf.Name, "", 0, nil)
    }

    // When an item is entered, verify if there is already a corresponding cluster page
    // If not, call NewPageCluster with the name of item and same data
    // If yes, switch to that cluster page
    list.SetSelectedFunc(func(i int, clusterName string, t string, s rune) {
        if data.App.GetPages().HasPage(clusterName) == false {
            // Attention: at this moment, data is the same
            // But after, if there is any specified fields related to each cluster
            // We have to create a new types.PageClusterData
            newData := &types.PageClusterData{Configs: data.Configs, App: data.App}
            newData.PageName = clusterName
            newPage, err := NewPageCluster(newData)
            if err != nil {
                fmt.Errorf("Cannot create new page %s\n", clusterName)
            }

            // Add new page
            data.App.GetPages().AddPage(clusterName, newPage, true, true)
            // Set focus in the new page to corresponding item
            l, err := newPage.GetElemList("list_clusters")
            if err != nil {
                fmt.Errorf("List list_cluster not found")
            }

            l.SetCurrentItem(i)
            //data.App.GetPages.SwitchToPage(clusterName)
            //data.App.SetFocus(newPage)
        } else {
            // If the page was already created, switch to it then
            data.App.GetPages().SwitchToPage(clusterName)
            // Set the corresponding item
            p, _ := data.App.GetPages().GetPage(clusterName)
            l, _ := p.GetElemList("list_clusters")
            l.SetCurrentItem(i)

        }
    })

    // Modify certain key events before forwarding others to default handler
    list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        //fmt.Println("Key pressed")
        switch event.Key() {
        case tcell.KeyTab:
            list, _ := page.GetElemList("list_menu")
            data.App.SetFocus(list)
            return nil
        }
        return event
    })

    return list
}

func createListMenu(page *console.Page) *tview.List {
    data := page.GetData().(*types.PageClusterData)
    clusterConfig := getClusterConfig(data.PageName, data)

    // Set column Clusters
    list := tview.NewList().ShowSecondaryText(false)
    list.SetBorder(true).SetTitle("Menu")
    //clusters.SetBackgroundColor(tcell.ColorDarkViolet)

    // Populate list content
    if clusterConfig.Type == "swarm" {
        // TODO: Must be filled up from menu registry of module swarm or k8s
        list.AddItem("services", "", 0, nil)
        list.AddItem("networks", "", 0, nil)
        list.AddItem("nodes", "", 0, nil)
    }

    // When clicking on an item
    list.SetSelectedFunc(func(i int, menuName string, t string, s rune) {
        if menuName == "services" {
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
        }
    })

    // Modify certain key events before forwarding others to default handler
    list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        //fmt.Println("Key pressed")
        switch event.Key() {
        case tcell.KeyTab:
            list, _ := page.GetElemList("list_clusters")
            data.App.SetFocus(list)
            return nil
        }
        return event
    })

    return list
}

func createTableService(page *console.Page) *tview.Table {
    //data := page.GetData().(*PageClusterData)
    //clusterConfig := getClusterConfig(data.PageName, data)

    // Set column Clusters
    table := tview.NewTable()
    table.SetBorders(true).SetTitle("Services")

    table.SetCell(0, 0, &tview.TableCell{Text: "ID", Align: tview.AlignCenter, Color: tcell.ColorYellow}).
          SetCell(0, 1, &tview.TableCell{Text: "Name", Align: tview.AlignCenter, Color: tcell.ColorYellow}).
          SetCell(0, 2, &tview.TableCell{Text: "Mode", Align: tview.AlignCenter, Color: tcell.ColorYellow}).
          SetCell(0, 3, &tview.TableCell{Text: "Replicas", Align: tview.AlignCenter, Color: tcell.ColorYellow}).
          SetCell(0, 4, &tview.TableCell{Text: "Image", Align: tview.AlignCenter, Color: tcell.ColorYellow})

    return table
}


//func initClusterClient(page *console.Page) error {
    //var client interface{}
    //var err error

    //data := page.GetData().(*types.PageClusterData)
    //config := getClusterConfig(data.PageName, data)

    //if config.Type == "docker" {
        //client, err = initSwarmClient(config)
        //if err != nil {
            //fmt.Println(err)
            //return err
        //}
    //}

    //data.Module.Client = client

    //return nil
//}

// initDockerClient initializes a docker client to remote cluster
// following authentication configuration
func initSwarmClient(config types.ClusterConfig) (interface{}, error) {
    var client interface{}
    //var err error

    if config.AuthType == "ssh" {
        if config.Auth.Type == "key" {
            sshConfig, err := ssh.NewClientConfigWithKey(config.Auth.Username, config.Auth.SshKey, "", false)
            if err != nil {
                fmt.Println(err)
                return nil, err
            }

            client, err = docker.NewSSHClient(config.Host, "/var/run/docker.sock", "1.30", sshConfig)
            if err != nil {
                fmt.Println(err)
                return nil, err
            }
        }
    }
    return client, nil
}

func getClusterConfig(name string, data *types.PageClusterData) types.ClusterConfig {
    for _, c := range data.Configs {
        if c.Name == name {
            return c
        }
    }

    return types.ClusterConfig{}
}


