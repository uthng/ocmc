package console

import (
    //"fmt"

    "github.com/gdamore/tcell"
    "github.com/rivo/tview"

)

type App struct {
    *tview.Application

    pages                               *Pages
    toolbar                             *Toolbar

    //console                             *tview.Application // The tview application.
    //pages                               *tview.Pages       // The application pages.
    //toolbar                             *console.Toolbar
    //finderFocus                 tview.Primitive    // The primitive in the Finder that last had focus.
}

func NewApp() *App {
    app := &App {
        Application: tview.NewApplication(),
        pages: NewPages(),
        toolbar: NewToolbar(),
    }

    flex := tview.NewFlex().SetDirection(tview.FlexRow)
    flex.AddItem(app.pages, 0, 15, false).
        AddItem(app.toolbar, 0, 1, false)

    initToolbar(app.toolbar)

    app.SetRoot(flex, true)

    // Override global Styles variable for look & feel
    tview.Styles = struct {
        PrimitiveBackgroundColor    tcell.Color // Main background color for primitives.
        ContrastBackgroundColor     tcell.Color // Background color for contrasting elements.
        MoreContrastBackgroundColor tcell.Color // Background color for even more contrasting elements.
        BorderColor                 tcell.Color // Box borders.
        TitleColor                  tcell.Color // Box titles.
        GraphicsColor               tcell.Color // Graphics.
        PrimaryTextColor            tcell.Color // Primary text.
        SecondaryTextColor          tcell.Color // Secondary text (e.g. labels).
        TertiaryTextColor           tcell.Color // Tertiary text (e.g. subtitles, notes).
        InverseTextColor            tcell.Color // Text on primary-colored backgrounds.
    }{
        PrimitiveBackgroundColor:    tcell.ColorBlack,
        ContrastBackgroundColor:     tcell.ColorBlue,
        MoreContrastBackgroundColor: tcell.ColorGreen,
        BorderColor:                 tcell.ColorWhite,
        TitleColor:                  tcell.ColorWhite,
        GraphicsColor:               tcell.ColorWhite,
        PrimaryTextColor:            tcell.ColorWhite,
        SecondaryTextColor:          tcell.ColorYellow,
        TertiaryTextColor:           tcell.ColorGreen,
        InverseTextColor:            tcell.ColorBlue,
    }
    return app
}

func (a *App) GetPages() *Pages {
    return a.pages
}

func (a *App) GetToolbar() *Toolbar {
    return a.toolbar
}

///////////////// PRIVATE FUNC /////////////////////////////////////////////
func initToolbar(toolbar *Toolbar) {
    btnHelp := tview.NewButton("[black]F1 [white]Help")
    toolbar.AddButton(btnHelp, 0)

    btnProject := tview.NewButton("[black]F2 [white]Projects")
    toolbar.AddButton(btnProject, 0)

    btnMenu := tview.NewButton("[black]F3 [white]Menu")
    toolbar.AddButton(btnMenu, 0)

    btnPagePrev := tview.NewButton("[black]PgDn [white]PgPrev")
    toolbar.AddButton(btnPagePrev, 0)

    btnPageNext := tview.NewButton("[black]PgUp [white]PgNext")
    toolbar.AddButton(btnPageNext, 0)

    btnQuit := tview.NewButton("[black]F10 [white]Quit")
    toolbar.AddButton(btnQuit, 0)

}
