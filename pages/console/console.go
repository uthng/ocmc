package console

import (
    "fmt"
    "strings"
    //"io"

    "github.com/rivo/tview"
    "github.com/gdamore/tcell"

    "github.com/docker/docker/client"

    "github.com/uthng/common/ssh"

    "github.com/uthng/ocmc/console"
    "github.com/uthng/ocmc/types"
    "github.com/uthng/ocmc/common/docker"
)

var selectedContainerId string
var sshClient           *ssh.Client

// NewPageConsole returns a new page console
func NewPageConsole(data *types.PageConsoleData) (*console.Page, error) {
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

    // Set column Containers
    err = setupListContainers("main", page)
    if err != nil {
        return nil, err
    }

    // Add container for console zone
    err = page.AddContainer("main", "console", tview.FlexRow, 0, 2, false)
    if err != nil {
        return nil, err
    }

    // Setup console zone
    err = setupContainerConsole("console", page)
    if err != nil {
        return nil, err
    }

    // Setup a connection client to server following authentication type
    config := data.Node.Config
    if config.Auth.Type == "ssh" {
        err = newSSHClient(config)
        if err != nil {
            return nil, err
        }
    }

    // Point to 1st element of container list
    list, _ := page.GetElemList("list_containers")
    list.SetCurrentItem(0)

    return page, nil
}


///////////// PRIVATE FUNCTIONS ///////////////////
func newSSHClient(config types.ConnConfig) error {
    if config.Auth.Kind == "key" {
        sshConfig, err := ssh.NewClientConfigWithKeyFile(config.Auth.Username, config.Auth.SshKey, config.Host, config.Port, false)
        if err != nil {
            return err
        }

        sshClient, err = ssh.NewClient(sshConfig)
        if err != nil {
            return err
        }

        return nil
    }

    return fmt.Errorf("Authentication type %s non supported\n", config.Auth.Kind)
}

// setupListContainers initializes and populates a list containing
// container's ID and name
func setupListContainers(container string, page *console.Page) error {
    var list *tview.List
    var mapContainers = make(map[string]string)

    data := page.GetData().(*types.PageConsoleData)

    // Check if list already exists. If not, create it. Otherwise reuse it
    list, err := page.GetElemList("list_containers")
    if err != nil {
        // Set column Clusters
        list = tview.NewList()
        list.SetBorder(true).SetBorderPadding(0, 0, 0, 0).SetTitle("Containers")
        list.ShowSecondaryText(true)

        err = page.AddItem(container, "list_containers", list, 0, 1, true)
        if err != nil {
            return err
        }
    }

    // Check if the node client is docker
    if data.Node.Config.Type == "docker" {
        // Get containers
        containers, err := docker.GetContainers(data.Node.Client.(*client.Client))
        if err != nil {
            return err
        }

        // Build map with ID as key and Name as value
        for _, c := range containers {
            mapContainers[c.ID] = strings.Trim(c.Names[0], "/")
        }

    }

    // Add containers to list
    for k, v := range mapContainers {
        list.AddItem(k[:12], v, 0, nil)
    }

    // Perform actions when an item is selected
    list.SetChangedFunc(func(i int, menuName string, t string, s rune) {
        // When a Enter occured on an item, clear input command
        // and textview output empty.
        inputField, _ := page.GetElemInputField("inputfield_command")
        textView, _ := page.GetElemTextView("textview_output")

        inputField.SetText("")
        textView.Clear()

        // Add container name to the title
        title, _ := page.GetElemTextView("title")
        title.Clear()
        fmt.Fprintf(title, "%s", data.PageName + " - " + t)

        // Set variable global
        selectedContainerId = menuName
    })

    // Modify certain key events before forwarding others to default handler
    list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        //fmt.Println("Key pressed")
        switch event.Key() {
        case tcell.KeyEsc:
            textView, _ := page.GetElemTextView("textview_output")
            data.App.SetFocus(textView)
            return nil
         case tcell.KeyTab:
            inputField, _ := page.GetElemInputField("inputfield_command")
            data.App.SetFocus(inputField)
            return nil

        }
        return event
    })

    return nil
}

func setupContainerConsole(container string, page *console.Page) error {
    err := setupInputFieldCommand(container, page)
    if err != nil {
        return err
    }

    err = setupTextViewOutput(container, page)
    if err != nil {
        return err
    }

    return nil
}

// setupInputFieldCommand setups an input field for command prompt
// and handles key events
func setupInputFieldCommand(container string, page *console.Page) error {
    var inputField *tview.InputField

    data := page.GetData().(*types.PageConsoleData)

    // Check if list already exists. If not, create it. Otherwise reuse it
    inputField, err := page.GetElemInputField("inputfield_command")
    if err != nil {
        // Set column Clusters
        inputField = tview.NewInputField()
        inputField.SetBorder(true).SetBorderPadding(0, 0, 0, 0)
        inputField.SetLabel("Command: ")

        err = page.AddItem(container, "inputfield_command", inputField, 0, 1, false)
        if err != nil {
            return err
        }
    }

    // Perform actions when Enter is pressed
    inputField.SetDoneFunc(func(key tcell.Key) {
        switch key {
        case tcell.KeyEnter:
            cmd := inputField.GetText()
            if len(cmd) > 0 {
                res, err := execCommand(selectedContainerId, cmd, data)
                // If no error, send command result to output textview
                if err != nil {
                    outputTextViewResponse(selectedContainerId, cmd, []byte(err.Error()), page)
                } else {
                    outputTextViewResponse(selectedContainerId, cmd, res, page)
                }
                // Clear input zone after entered
                inputField.SetText("")
            }
        }
    })

    // Modify certain key events before forwarding others to default handler
    inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        //fmt.Println("Key pressed")
        switch event.Key() {
        case tcell.KeyEsc:
            list, _ := page.GetElemList("list_containers")
            data.App.SetFocus(list)
            return nil
         case tcell.KeyTab:
            textView, _ := page.GetElemTextView("textview_output")
            data.App.SetFocus(textView)
            return nil

        }
        return event
    })


    return nil
}

// setupTextViewOutput initializes a text zone for command output
// and handles key events
func setupTextViewOutput(container string, page *console.Page) error {
    var textView *tview.TextView

    data := page.GetData().(*types.PageConsoleData)

    // Check if list already exists. If not, create it. Otherwise reuse it
    textView, err := page.GetElemTextView("textview_output")
    if err != nil {
        // Set column Clusters
        textView = tview.NewTextView()
        textView.SetBorder(true).SetBorderPadding(0, 0, 0, 0)

        err = page.AddItem(container, "textview_output", textView, 0, 10, false)
        if err != nil {
            return err
        }
    }

    // Modify certain key events before forwarding others to default handler
   textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        //fmt.Println("Key pressed")
        switch event.Key() {
        case tcell.KeyEsc:
            inputField, _ := page.GetElemInputField("inputfield_command")
            data.App.SetFocus(inputField)
            return nil
         case tcell.KeyTab:
            list, _ := page.GetElemList("list_containers")
            data.App.SetFocus(list)
            return nil

        }
        return event
    })

    return nil
}

func outputTextViewResponse(cid string, cmd string, response []byte, page *console.Page) {
    output, _ := page.GetElemTextView("textview_output")
    output.Write([]byte(cid[:12] + " $ " + cmd + "\n"))
    output.Write(response)
    output.Write([]byte("\n\n"))
}
