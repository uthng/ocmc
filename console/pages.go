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

type Pages struct {
    *tview.Pages

    CurrentPage         string
    //pages               map[string]*Page
    pages               *list.List
}

func NewPages() *Pages {
    p := &Pages {
        Pages: tview.NewPages(),
    }

    //p.pages = make(map[string]*Page)
    p.pages = list.New()

    return p
}

func (p *Pages) AddPage(name string, page *Page, resize bool, visible bool) {
    p.Pages.AddPage(name, page, resize, visible)

    //p.pages[name] = page
    p.pages.PushBack(page)
    p.CurrentPage = name
}

func (p *Pages) GetPage(name string) (*Page, error) {
    err := p.Pages.HasPage(name)
    if err == false {
        return nil, ErrPageNotFound
    }

    //page, ok := p.pages[name]
    //if ok == false || err == false {
        //return nil, ErrPageNotFound
    //}

    // Iterate through list and print its contents.
    for e := p.pages.Front(); e != nil; e = e.Next() {
        page := e.Value.(*Page)
         if page.Name == name {
            return page, nil
         }
    }

    return nil, ErrPageNotFound
}

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

