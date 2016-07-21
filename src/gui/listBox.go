package gui

import (
	"log"
	"strconv"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type ListBox struct {
	treeView      *gtk.TreeView
	listStore     *gtk.ListStore
	selIndex      int
	funOnClick    func(ind int)
	funOnDblClick func(ind int)
}

func NewListBox() *ListBox {
	lb := &ListBox{}
	render, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Println("Error create CellRendererTextNew", err)
		return nil
	}
	columns, err := gtk.TreeViewColumnNewWithAttribute("Radio stations", render, "text", 0)
	if err != nil {
		log.Println("Error create TreeViewColumnNewWithAttribute", err)
		return nil
	}
	lb.treeView, err = gtk.TreeViewNew()
	if err != nil {
		log.Println("Error create TreeViewNew", err)
		return nil
	}
	lb.treeView.AppendColumn(columns)

	lb.listStore, err = gtk.ListStoreNew(glib.TYPE_STRING)
	if err != nil {
		log.Fatal("Error create ListStoreNew", err)
	}
	lb.treeView.SetModel(lb.listStore)

	lb.treeView.Connect("cursor-changed", func() {
		if lb.funOnClick != nil {
			lb.funOnClick(lb.GetSelected())
		}
	})

	lb.treeView.Connect("row_activated", func() {
		if lb.funOnDblClick != nil {
			lb.funOnDblClick(lb.GetSelected())
		}
	})

	return lb
}

func (l *ListBox) GetSelected() int {
	path, _ := l.treeView.GetCursor()
	if path != nil {
		inds := path.GetIndices()
		if len(inds) > 0 {
			return inds[0]
		}
	}
	return -1
}

func (l *ListBox) SetSelected(ind int) {
	if ind == -1 {
		return
	}
	path, err := gtk.TreePathNewFromString(strconv.Itoa(ind))
	if err == nil {
		l.treeView.SetCursor(path, nil, false)
	}
}

func (l *ListBox) Update(list []string) {
	l.listStore.Clear()
	for _, s := range list {
		l.listStore.SetValue(l.listStore.Append(), 0, s)
	}
}

func (l *ListBox) OnClick(cb func(ind int)) {
	l.funOnClick = cb
}

func (l *ListBox) OnDblClick(cb func(ind int)) {
	l.funOnDblClick = cb
}

func (l *ListBox) GetWidget() gtk.IWidget {
	return l.treeView
}
