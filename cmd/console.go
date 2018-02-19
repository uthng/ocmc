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
    //"github.com/uthng/ocmc/types"
    "github.com/uthng/ocmc/common/config"

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
    app.GetPages().AddPage("home", pageHome, true, true)

    // Build 1st page cluster
    data := config.ReadClusterConfigFromFile()
    data.App = app
    data.PageName = data.Configs[0].Name
    if len(data.Configs) > 0 {
        pageCluster, err := pages.NewPageCluster(data)
        if err != nil {
            panic(err)
        }
        app.GetPages().AddPage(data.Configs[0].Name, pageCluster, true, true)
    }

    //app.GetPages().SwitchToPage("home")

    if err := app.SetFocus(app.GetPages()).Run(); err != nil {
        panic("Error running application")
    }
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
        //app.GetPages().SwitchToPage("cluster")
        //page, _ := app.GetPages().GetPage("TVL PPD FR")
        //clusters, _ := page.GetElemList("list_clusters")

        //app.SetFocus(clusters)

    case tcell.KeyF10:
        app.Stop()
    }

    return event
}
