package console

import (
    "fmt"
    "errors"
    "reflect"

    //"github.com/gdamore/tcell"
    "github.com/rivo/tview"
)

var (
    ErrContainerNotFound = errors.New("Container not found")
    ErrElemNotFound = errors.New("Element not found")
    ErrElemIncorrectType = errors.New("Element incorrect type")
)

type Page struct {
    *tview.Flex

    containers map[string]*tview.Flex
    elems map[string]interface{}
}

func NewPage() *Page {
    p := &Page {
        Flex: tview.NewFlex(),
    }

    p.containers = make(map[string]*tview.Flex)
    p.elems = make(map[string]interface{})

    return p
}

func (p *Page) AddContainer(parent string, name string, direction int, size int, proportion int, focus bool) error {
    var flex *tview.Flex

    flex = tview.NewFlex().SetDirection(direction)
    if parent == "root" || parent == "" {
        p.Flex.AddItem(flex, size, proportion, focus)
    } else {
        parentFlex, ok := p.containers[parent]
        if ok == false {
            return ErrContainerNotFound
        }

        parentFlex.AddItem(flex, size, proportion, focus)
    }

    // Add new container to container list
    p.containers[name] = flex

    return nil
}

func (p *Page) AddItem(container string, name string, item interface{}, size int, proportion int, focus bool) error {
    var primitive tview.Primitive

    switch item.(type) {
    case *tview.Box:
        primitive = item.(*tview.Box)
    case *tview.Button:
        primitive = item.(*tview.Button)
    case *tview.Checkbox:
        primitive = item.(*tview.Checkbox)
    case *tview.DropDown:
        primitive = item.(*tview.DropDown)
    case *tview.Flex:
        primitive = item.(*tview.Flex)
    case *tview.Form:
        primitive = item.(*tview.Form)
    case *tview.Frame:
        primitive = item.(*tview.Frame)
    case *tview.InputField:
        primitive = item.(*tview.InputField)
    case *tview.List:
        primitive = item.(*tview.List)
    case *tview.Modal:
        primitive = item.(*tview.Modal)
    case *tview.Pages:
        primitive = item.(*tview.Pages)
    case *tview.Table:
        primitive = item.(*tview.Table)
    case *tview.TextView:
        primitive = item.(*tview.TextView)
    default:
        return ErrElemIncorrectType
    }

    if container == "root" || container == "" {
        p.Flex.AddItem(primitive, size, proportion, focus)
    } else {
        flex, ok := p.containers[container]
        if ok == false {
            fmt.Printf("%v\n", p.containers)
            return fmt.Errorf("Container %s not found", container)
        }

        flex.AddItem(primitive, size, proportion, focus)
    }

    p.elems[name] = item

    return nil
}

func (p *Page) GetContainer(name string) (*tview.Flex, error) {
    container, ok := p.containers[name]
    if ok == false {
        return nil, ErrContainerNotFound
    }

    return container, nil
}

func (p *Page) SetContainerDirection(container string, direction int) error {
    if container == "root" || container == "" {
        p.SetDirection(direction)
    } else {
        flex, ok := p.containers[container]
        if ok == false {
            return ErrContainerNotFound
        }

        flex.SetDirection(direction)
    }

    return nil
}

func (p *Page) GetElemBox(name string) (*tview.Box, error) {
    elem, err := getElem(name, p.elems)
    if err != nil {
        return nil, err
    }

    typeElem := reflect.TypeOf(elem).String()
    if typeElem != "*tview.Box" {
        return nil, ErrElemIncorrectType
    }

    return elem.(*tview.Box), nil
}

func (p *Page) GetElemButton(name string) (*tview.Button, error) {
    elem, err := getElem(name, p.elems)
    if err != nil {
        return nil, err
    }

    typeElem := reflect.TypeOf(elem).String()
    if typeElem != "*tview.Button" {
        return nil, ErrElemIncorrectType
    }

    return elem.(*tview.Button), nil
}

func (p *Page) GetElemCheckbox(name string) (*tview.Checkbox, error) {
    elem, err := getElem(name, p.elems)
    if err != nil {
        return nil, err
    }

    typeElem := reflect.TypeOf(elem).String()
    if typeElem != "*tview.Checkbox" {
        return nil, ErrElemIncorrectType
    }

    return elem.(*tview.Checkbox), nil
}

func (p *Page) GetElemDropDown(name string) (*tview.DropDown, error) {
    elem, err := getElem(name, p.elems)
    if err != nil {
        return nil, err
    }

    typeElem := reflect.TypeOf(elem).String()
    if typeElem != "*tview.DropDown" {
        return nil, ErrElemIncorrectType
    }

    return elem.(*tview.DropDown), nil
}

func (p *Page) GetElemFlex(name string) (*tview.Flex, error) {
    elem, err := getElem(name, p.elems)
    if err != nil {
        return nil, err
    }

    typeElem := reflect.TypeOf(elem).String()
    if typeElem != "*tview.Flex" {
        return nil, ErrElemIncorrectType
    }

    return elem.(*tview.Flex), nil
}

func (p *Page) GetElemForm(name string) (*tview.Form, error) {
    elem, err := getElem(name, p.elems)
    if err != nil {
        return nil, err
    }

    typeElem := reflect.TypeOf(elem).String()
    if typeElem != "*tview.Form" {
        return nil, ErrElemIncorrectType
    }

    return elem.(*tview.Form), nil
}

func (p *Page) GetElemFrame(name string) (*tview.Frame, error) {
    elem, err := getElem(name, p.elems)
    if err != nil {
        return nil, err
    }

    typeElem := reflect.TypeOf(elem).String()
    if typeElem != "*tview.Frame" {
        return nil, ErrElemIncorrectType
    }

    return elem.(*tview.Frame), nil
}

func (p *Page) GetElemInputField(name string) (*tview.InputField, error) {
    elem, err := getElem(name, p.elems)
    if err != nil {
        return nil, err
    }

    typeElem := reflect.TypeOf(elem).String()
    if typeElem != "*tview.InputField" {
        return nil, ErrElemIncorrectType
    }

    return elem.(*tview.InputField), nil
}

func (p *Page) GetElemList(name string) (*tview.List, error) {
    elem, err := getElem(name, p.elems)
    if err != nil {
        return nil, err
    }

    typeElem := reflect.TypeOf(elem).String()
    if typeElem != "*tview.List" {
        return nil, ErrElemIncorrectType
    }

    return elem.(*tview.List), nil
}

func (p *Page) GetElemModal(name string) (*tview.Modal, error) {
    elem, err := getElem(name, p.elems)
    if err != nil {
        return nil, err
    }

    typeElem := reflect.TypeOf(elem).String()
    if typeElem != "*tview.Modal" {
        return nil, ErrElemIncorrectType
    }

    return elem.(*tview.Modal), nil
}

func (p *Page) GetElemPages(name string) (*tview.Pages, error) {
    elem, err := getElem(name, p.elems)
    if err != nil {
        return nil, err
    }

    typeElem := reflect.TypeOf(elem).String()
    if typeElem != "*tview.Pages" {
        return nil, ErrElemIncorrectType
    }

    return elem.(*tview.Pages), nil
}

func (p *Page) GetElemTable(name string) (*tview.Table, error) {
    elem, err := getElem(name, p.elems)
    if err != nil {
        return nil, err
    }

    typeElem := reflect.TypeOf(elem).String()
    if typeElem != "*tview.Table" {
        return nil, ErrElemIncorrectType
    }

    return elem.(*tview.Table), nil
}

func (p *Page) GetElemTableCell(name string) (*tview.TableCell, error) {
    elem, err := getElem(name, p.elems)
    if err != nil {
        return nil, err
    }

    typeElem := reflect.TypeOf(elem).String()
    if typeElem != "*tview.TableCell" {
        return nil, ErrElemIncorrectType
    }

    return elem.(*tview.TableCell), nil
}

func (p *Page) GetElemTextView(name string) (*tview.TextView, error) {
    elem, err := getElem(name, p.elems)
    if err != nil {
        return nil, err
    }

    typeElem := reflect.TypeOf(elem).String()
    if typeElem != "*tview.TextView" {
        return nil, ErrElemIncorrectType
    }

    return elem.(*tview.TextView), nil
}


func getElem(name string, elems map[string]interface{}) (interface{}, error) {
    elem, ok := elems[name]
    if ok == false {
        return nil, ErrElemNotFound
    }

    return elem, nil
}
