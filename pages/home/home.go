package home

import (
    "fmt"

    "github.com/rivo/tview"
    "github.com/gdamore/tcell"

    "github.com/uthng/ocmc/console"
)
const (
    textLogo = `
                                                                                            
                                                                                            
     OOOOOOOOO             CCCCCCCCCCCCCMMMMMMMM               MMMMMMMM        CCCCCCCCCCCCC
   OO:::::::::OO        CCC::::::::::::CM:::::::M             M:::::::M     CCC::::::::::::C
 OO:::::::::::::OO    CC:::::::::::::::CM::::::::M           M::::::::M   CC:::::::::::::::C
O:::::::OOO:::::::O  C:::::CCCCCCCC::::CM:::::::::M         M:::::::::M  C:::::CCCCCCCC::::C
O::::::O   O::::::O C:::::C       CCCCCCM::::::::::M       M::::::::::M C:::::C       CCCCCC
O:::::O     O:::::OC:::::C              M:::::::::::M     M:::::::::::MC:::::C              
O:::::O     O:::::OC:::::C              M:::::::M::::M   M::::M:::::::MC:::::C              
O:::::O     O:::::OC:::::C              M::::::M M::::M M::::M M::::::MC:::::C              
O:::::O     O:::::OC:::::C              M::::::M  M::::M::::M  M::::::MC:::::C              
O:::::O     O:::::OC:::::C              M::::::M   M:::::::M   M::::::MC:::::C              
O:::::O     O:::::OC:::::C              M::::::M    M:::::M    M::::::MC:::::C              
O::::::O   O::::::O C:::::C       CCCCCCM::::::M     MMMMM     M::::::M C:::::C       CCCCCC
O:::::::OOO:::::::O  C:::::CCCCCCCC::::CM::::::M               M::::::M  C:::::CCCCCCCC::::C
 OO:::::::::::::OO    CC:::::::::::::::CM::::::M               M::::::M   CC:::::::::::::::C
   OO:::::::::OO        CCC::::::::::::CM::::::M               M::::::M     CCC::::::::::::C
     OOOOOOOOO             CCCCCCCCCCCCCMMMMMMMM               MMMMMMMM        CCCCCCCCCCCCC
                                                                                            
                                                                                            
                                                                                            
                                                                                            
                                                                                            
                                                                                            
                                                                                            
+-++-++-++-++-++-++-++-++-++-++-++-+ +-++-++-++-++-++-++-+ +-++-++-++-++-++-++-++-++-++-+ +-++-++-++-++-++-++-+
|O||r||c||h||e||s||t||r||a||t||o||r| |C||l||u||s||t||e||r| |M||a||n||a||g||e||m||e||n||t| |C||o||n||s||o||l||e|
+-++-++-++-++-++-++-++-++-++-++-++-+ +-++-++-++-++-++-++-+ +-++-++-++-++-++-++-++-++-++-+ +-++-++-++-++-++-++-+
     `
)

func NewPageHome(name string) (*console.Page, error) {
    page := console.NewPage(name)

    textView := tview.NewTextView().
                      SetDynamicColors(true).
                      SetRegions(true)

    fmt.Fprintf(textView, "%s", textLogo)
    textView.SetBorder(true)
    textView.SetTextAlign(tview.AlignCenter).SetTextColor(tcell.ColorDarkViolet)

    err := page.AddItem("root", "logo", textView, 0, 1, false)
    if err != nil {
        return nil, err
    }

    return page, nil
}
