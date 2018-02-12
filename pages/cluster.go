package pages

import (
    "fmt"

    "github.com/rivo/tview"
    "github.com/gdamore/tcell"

    "golang.org/x/net/context"

    // docker
    "github.com/moby/moby/client"
    "github.com/docker/docker/api/types"
    //"github.com/docker/docker/api/types/filters"
    //"github.com/docker/docker/api/types/swarm"

    "github.com/uthng/ocmc/console"
    "github.com/uthng/common/ssh"
    "github.com/uthng/common/docker"
)

type ClusterConfig struct {
    Name        string
    Host        string
    Type        string
    AuthType    string
    Auth        ClusterAuth
}

type ClusterAuth struct {
    Type        string
    Ca          string
    Client      string
    ClientKey   string

    SshKey      string

    Username    string
    Passord     string
}

type PageClusterData struct {
    PageName            string
    Configs             []ClusterConfig

    Client              interface{}

    App                 *console.App
}

func NewPageCluster(data *PageClusterData) (*console.Page, error) {
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

    // Setup menu for project
    listMenu := createListMenu(page)

    err = page.AddItem("cluster", "list_menu", listMenu, 0, 1, false)
    if err != nil {
        return nil, err
    }

    // Setup zone showing project details
    //details := tview.NewBox()
    //details.SetBorder(true).SetTitle("Details")

    //err = page.AddItem("cluster", "details", details, 0, 5, false)
    //if err != nil {
        //return nil, err
    //}

    // Add container for cluster
    err = page.AddContainer("cluster", "details", tview.FlexColumn, 0, 10, false)
    if err != nil {
        return nil, err
    }

    

    // Init cluster client
    initClusterClient(page)
    return page, nil
}


///////////// PRIVATE FUNCTIONS ///////////////////
func createListCluster(page *console.Page) *tview.List {
    data := page.GetData().(*PageClusterData)

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
            // We have to create a new PageClusterData
            newData := &PageClusterData{Configs: data.Configs, App: data.App}
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
        case tcell.KeyRight:
            list, _ := page.GetElemList("list_menu")
            data.App.SetFocus(list)
            return nil
        }
        return event
    })

    return list
}

func createListMenu(page *console.Page) *tview.List {
    data := page.GetData().(*PageClusterData)
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
            client := data.Client.(*client.Client)
            services, err := client.ServiceList(context.Background(), types.ServiceListOptions{})
            if err != nil {
                fmt.Println(err)
            }

            for _, service := range services {
                fmt.Println(service.Spec.Annotations.Name)
            }

        }
    })

    // Modify certain key events before forwarding others to default handler
    list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        //fmt.Println("Key pressed")
        switch event.Key() {
        case tcell.KeyLeft:
            list, _ := page.GetElemList("list_clusters")
            data.App.SetFocus(list)
            return nil
        }
        return event
    })

    return list
}

func initClusterClient(page *console.Page) error {
    var client interface{}
    var err error

    data := page.GetData().(*PageClusterData)
    config := getClusterConfig(data.PageName, data)

    if config.Type == "swarm" {
        client, err = initSwarmClient(config)
        if err != nil {
            fmt.Println(err)
            return err
        }
    }

    data.Client = client

    return nil
}

// initDockerClient initializes a docker client to remote cluster
// following authentication configuration
func initSwarmClient(config ClusterConfig) (interface{}, error) {
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

func getClusterConfig(name string, data *PageClusterData) ClusterConfig {
    for _, c := range data.Configs {
        if c.Name == name {
            return c
        }
    }

    return ClusterConfig{}
}


