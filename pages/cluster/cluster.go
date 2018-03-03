package cluster

import (
    "fmt"
    //"strings"

    "github.com/rivo/tview"
    "github.com/gdamore/tcell"

    "github.com/uthng/ocmc/console"
    "github.com/uthng/ocmc/types"
    "github.com/uthng/ocmc/common/config"
    module_docker "github.com/uthng/ocmc/modules/docker"
    module_k8s "github.com/uthng/ocmc/modules/k8s"
)


func NewPageCluster(data *types.PageClusterData) (*console.Page, error) {
    page := console.NewPage(data.PageName)

    // Set page data
    page.Data = data

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
        data.Module, err = module_docker.NewModuleDocker(connConfig)
    } else if connConfig.Type == "k8s" {
        data.Module, err = module_k8s.NewModuleK8s(connConfig)
    }
    if err != nil {
        return nil, err
    }

    data.Module.Layout("cluster", page)

    return page, nil
}


///////////// PRIVATE FUNCTIONS ///////////////////
func createListCluster(page *console.Page) *tview.List {
    data := page.Data.(*types.PageClusterData)

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
                fmt.Println(err)
            } else {
                // Add new page
                data.App.GetPages().AddPage(clusterName, newPage, true, true)
                // Set focus in the new page to corresponding item
                l, err := newPage.GetElemList("list_clusters")
                if err != nil {
                    fmt.Println(err)
                }

                l.SetCurrentItem(i)
                //data.App.GetPages.SwitchToPage(clusterName)
                //data.App.SetFocus(newPage)
            }
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
        switch event.Key() {
        case tcell.KeyTab:
            list, _ := page.GetElemList("list_menu")
            data.App.SetFocus(list)
            return nil
        case tcell.KeyPgUp, tcell.KeyPgDn:
            // Loop list of configs to find out the index of the current page
            // because it is the same for list_clusters
            for i, conf := range data.Configs {
                if conf.Name == data.App.GetPages().CurrentPage {
                    // Get the page corresponding to the current page
                    p, _ := data.App.GetPages().GetPage(data.App.GetPages().CurrentPage)
                    // Get list cluster of current page
                    l, _ := p.GetElemList("list_clusters")
                    // Set the item corresponding to the current page
                    l.SetCurrentItem(i)
                }
            }
            return nil
        }
        return event
    })

    return list
}
