package pages

import (
    "fmt"

    "github.com/rivo/tview"
    //"github.com/gdamore/tcell"

    "github.com/uthng/ocmc/console"
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

    Username    string
    Passord     string
}

type PageClusterData struct {

    Configs      []ClusterConfig

    //Connection
    App          *console.App
}

func NewPageCluster(name string, data PageClusterData) (*console.Page, error) {
    page := console.NewPage()

    // Set page data
    page.BindData(data)

    // Setup Title
    title := tview.NewTextView()
    fmt.Fprintf(title, "%s", name)
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
    listClusters := createListCluster(data)
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
    listMenu := tview.NewList().ShowSecondaryText(false)
    listMenu.SetBorder(true).SetTitle("Menu")

    err = page.AddItem("cluster", "list_menu", listMenu, 0, 1, false)
    if err != nil {
        return nil, err
    }

    // Setup zone showing project details
    details := tview.NewBox()
    details.SetBorder(true).SetTitle("Details")

    err = page.AddItem("cluster", "details", details, 0, 5, false)
    if err != nil {
        return nil, err
    }

    return page, nil
}


///////////// PRIVATE FUNCTIONS ///////////////////
func createListCluster(data PageClusterData) *tview.List {
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
            newData := PageClusterData{Configs: data.Configs, App: data.App}
            newPage, err := NewPageCluster(clusterName, newData)
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
            l, err := p.GetElemList("list_clusters")
            if err != nil {
                fmt.Errorf("List list_cluster not found")
            }

            l.SetCurrentItem(i)

        }
    })

    return list
}

//func createListMenu(data 
