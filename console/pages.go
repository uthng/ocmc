package console

import (
    //"fmt"
    "errors"

    //"github.com/gdamore/tcell"
    "github.com/rivo/tview"
)

var (
    ErrPageNotFound = errors.New("Page not found")
)

type Pages struct {
    *tview.Pages

    pages      map[string]*Page
}

func NewPages() *Pages {
    p := &Pages {
        Pages: tview.NewPages(),
    }

    p.pages = make(map[string]*Page)

    return p
}

func (p *Pages) AddPage(name string, page *Page, resize bool, visible bool) {
    p.Pages.AddPage(name, page, resize, visible)

    p.pages[name] = page
}

func (p *Pages) GetPage(name string) (*Page, error) {
    err := p.Pages.HasPage(name)

    page, ok := p.pages[name]
    if ok == false || err == false {
        return nil, ErrPageNotFound
    }

    return page, nil
}

func (p *Pages) RemovePage(name string) {
    p.Pages.RemovePage(name)

    delete(p.pages, name)
}
