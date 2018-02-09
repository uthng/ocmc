// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
    //"fmt"

    "github.com/spf13/cobra"

    "github.com/gdamore/tcell"

    "github.com/uthng/ocmc/console"
    "github.com/uthng/ocmc/pages"

)

var app *console.App

// consoleCmd represents the console command
var consoleCmd = &cobra.Command{
    Use:   "console",
    Short: "Orchestrator Cluster Management Console",
    Long: `OCMC is an application console to manage different orchestrator cluster
such as Kubernetes or Swarm`,
    Run: func(cmd *cobra.Command, args []string) {
        initApp()
    },
}

func init() {
    rootCmd.AddCommand(consoleCmd)

    // Here you will define your flags and configuration settings.

    // Cobra supports Persistent Flags which will work for this command
    // and all subcommands, e.g.:
    // consoleCmd.PersistentFlags().String("foo", "", "A help for foo")

    // Cobra supports local flags which will only run when this command
    // is called directly, e.g.:
    // consoleCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initApp() {
    app = console.NewApp()

    // Set handler to handle key event
    app.SetInputCapture(handlerKeyEvent)

    pageHome, err := pages.NewPageHome("Home")
    if err != nil {
        panic(err)
    }

    pageProject, err := pages.NewPageProject("Projects")
    if err != nil {
        panic(err)
    }

    app.GetPages().AddPage("home", pageHome, true, true)
    app.GetPages().AddPage("project", pageProject, true, true)

    app.GetPages().SwitchToPage("home")

    populateProjectList()
    //pageHome := createHomePage()
    //app.pages.AddPage("Home", pageHome, true, true)

    // Create Project Page
    //createProjectPage()

    // Initialize toolbar with app buttons
    //createToolbar()

    //app.pages.SwitchToPage("Home")

    // Run
    //app.console.SetRoot(flex, true)
    if err := app.Run(); err != nil {
        panic("Error running application")
    }
}

// createPage Setup all layout's elements
//func createHomePage() tview.Primitive {
    //textView := tview.NewTextView().
                //SetDynamicColors(true).
                //SetRegions(true)

    //fmt.Fprintf(textView, "%s", textPageHome)
    //textView.SetBorder(true)
    //textView.SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorDarkViolet)

    //return textView
//}

func populateProjectList() error {
    page, err := app.GetPages().GetPage("project")
    if err != nil {
        return err
    }

    projects, err := page.GetElemList("projects")
    if err != nil {
        return err
    }

    projects.AddItem("TVL PPD FR", "Travel PPD FR cluster", 0, nil)
    projects.AddItem("TVL PPD NL", "Travel PPD NL cluster", 0, nil)
    projects.AddItem("TVL PRD FR", "Travel PRD FR cluster", 0, nil)
    projects.AddItem("TVL PRD NL", "Travel PRD NL cluster", 0, nil)

    return nil
}

func handlerKeyEvent(event *tcell.EventKey) *tcell.EventKey {
    switch event.Key() {
    case tcell.KeyEscape:
        // Go back to Finder.
        //pageProject := console.NewProjectPage()
        //pageProject.Initialize("TOTO")
        //app.pages.AddPage("Project", pageProject, true, true)

        //app.pages.SwitchToPage("Project")
        //pages.SwitchToPage(finderPage)
        //if finderFocus != nil {
            //app.SetFocus(finderFocus)
        //}
    //case tcell.KeyEnter:
        //// Load the next batch of rows.
        //loadRows(table.GetRowCount() - 1)
        //table.ScrollToEnd()
        //app.toolbar.SetBorder(true)
    //}
    //case tcell.KeyPgUp:
    //case tcell.KeyPgDn:
    case tcell.KeyF2:
        app.GetPages().SwitchToPage("project")
        page, _ := app.GetPages().GetPage("project")
        projects, _ := page.GetElemList("projects")

        app.SetFocus(projects)

    case tcell.KeyF10:
        app.Stop()
    }

    return event
}
