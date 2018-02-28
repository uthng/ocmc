package console

import (
    //"fmt"
    "errors"
    "container/list"

    //"github.com/gdamore/tcell"
    "github.com/rivo/tview"
)

var (
    ErrPageNotFound = errors.New("Page not found")
)

// Page is a embedded type of tview.Page with 2 more fields:
// the current selected page and a page list
type Pages struct {
    *tview.Pages
    // Current selected page
    CurrentPage         string
    // page list
    pages               *list.List
}

// NewPages returns a new Pages
func NewPages() *Pages {
    p := &Pages {
        Pages: tview.NewPages(),
    }

    //p.pages = make(map[string]*Page)
    p.pages = list.New()

    return p
}

// AddPage adds a new page to tview.Pages and page list
func (p *Pages) AddPage(name string, page *Page, resize bool, visible bool) {
    p.Pages.AddPage(name, page, resize, visible)

    //p.pages[name] = page
    p.pages.PushBack(page)
    p.CurrentPage = name
}

// GetPage returns the page corresponding to the given name
// if the page is in tview.Pages slice and page list
func (p *Pages) GetPage(name string) (*Page, error) {
    err := p.Pages.HasPage(name)
    if err == false {
        return nil, ErrPageNotFound
    }

    // Iterate through list and print its contents.
    for e := p.pages.Front(); e != nil; e = e.Next() {
        page := e.Value.(*Page)
         if page.Name == name {
            return page, nil
         }
    }

    return nil, ErrPageNotFound
}

// RemovePage deletes the page corresponding to the given name
// if it exists
func (p *Pages) RemovePage(name string) {
    p.Pages.RemovePage(name)

    //delete(p.pages, name)
    for e := p.pages.Front(); e != nil; e = e.Next() {
        page := e.Value.(*Page)
         if page.Name == name {
            // Check if name is current page.
            // If yes, set current page to the next one
            // or previous one
            if p.CurrentPage == name {
                if e.Next() != nil {
                    p.SwitchToNextPage()
                } else if e.Prev() != nil {
                    p.SwitchToPrevPage()
                }
            }
            p.pages.Remove(e)
            return
         }
    }
}

// SwitchToNextPage switchs to the next page of the current page
func (p *Pages) SwitchToNextPage() {
    for e := p.pages.Front(); e != nil; e = e.Next() {
        page := e.Value.(*Page)
        if page.Name == p.CurrentPage {
            if e.Next() != nil {
                p.CurrentPage = e.Next().Value.(*Page).Name
                p.Pages.SwitchToPage(p.CurrentPage)
                return
            }
        }
    }
}

// SwitchToPrevPage switchs to the prev page of the current page
func (p *Pages) SwitchToPrevPage() {
    for e := p.pages.Front(); e != nil; e = e.Next() {
        page := e.Value.(*Page)
        if page.Name == p.CurrentPage {
            if e.Prev() != nil {
                p.CurrentPage = e.Prev().Value.(*Page).Name
                p.Pages.SwitchToPage(p.CurrentPage)
                return
            }
        }
    }
}

// SwitchToPage switchs to the page corresponding to the given name
func (p *Pages) SwitchToPage(name string) {
    for e := p.pages.Front(); e != nil; e = e.Next() {
        page := e.Value.(*Page)
        if page.Name == name {
            p.CurrentPage = name
            p.Pages.SwitchToPage(p.CurrentPage)
        }
    }
}

