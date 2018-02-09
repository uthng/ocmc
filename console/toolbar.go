package console

import (
    //"fmt"

    "github.com/rivo/tview"
)

//type ZoneProject interface {

//}

type Toolbar struct {
    *tview.Flex

    //buttons []*tview.Button
}

func NewToolbar() *Toolbar {
    toolbar := &Toolbar {
        Flex: tview.NewFlex(),
    }

    return toolbar
}

//func (toolbar *Toolbar) Initialize() {
    //var button *tview.Button

    //button = tview.NewButton("Close")


    //toolbar.Flex.AddItem(button, 10, 1, false)
    //toolbar.Flex.SetBorder(true)
//}

func (toolbar *Toolbar) AddButton(button *tview.Button, size int) {
    toolbar.Flex.AddItem(button, size, 1, false)
}
