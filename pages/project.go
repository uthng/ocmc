package pages

import (
    "fmt"

    "github.com/rivo/tview"
    //"github.com/gdamore/tcell"

    "github.com/uthng/ocmc/console"
)

func NewPageProject(name string) (*console.Page, error) {
    page := console.NewPage()

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
    err = page.AddContainer("root", "main", tview.FlexColumn, 0, 10, false)
    if err != nil {
        return nil, err
    }

    // Set column Projects
    projects := tview.NewList().ShowSecondaryText(false)
    projects.SetBorder(true).SetTitle("Projects")
    //projects.SetBackgroundColor(tcell.ColorDarkViolet)

    err = page.AddItem("main", "projects", projects, 0, 1, false)
    if err != nil {
        return nil, err
    }

    // Add container for project
    err = page.AddContainer("main", "project", tview.FlexColumn, 0, 7, false)
    if err != nil {
        return nil, err
    }

    // Setup menu for project
    menu := tview.NewList().ShowSecondaryText(false)
    menu.SetBorder(true).SetTitle("Menu")

    err = page.AddItem("project", "menu", menu, 0, 1, false)
    if err != nil {
        return nil, err
    }

    // Setup zone showing project details
    details := tview.NewBox()
    details.SetBorder(true).SetTitle("Details")

    err = page.AddItem("project", "details", details, 0, 5, false)
    if err != nil {
        return nil, err
    }

    return page, nil
}
