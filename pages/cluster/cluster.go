package cluster

import (
    "fmt"
    //"strings"

    "github.com/rivo/tview"
    "github.com/gdamore/tcell"

    "github.com/uthng/ocmc/console"
    "github.com/uthng/ocmc/types"
    "github.com/uthng/ocmc/common/config"
    "github.com/uthng/ocmc/common/docker"
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

    connConfig := config.GetClusterConfig(data.PageName, data).Config
    if connConfig.Type == "docker" {
        data.Module = module_docker.NewModuleDocker()
        data.Module.Client, err = docker.NewDockerClient(connConfig)
        if err != nil {
            fmt.Println(err)
            return nil, err
        }

        data.Module.Layout("cluster", page)
        //fmt.Println("data %v\n", data.Module)
        //fmt.Println("data %v\n", page.GetData().(*types.PageClusterData).Module)

    }

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
