package gui

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

func NewHBox(spacing int) *gtk.Box {
	box, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, spacing)
	if err != nil {
		log.Println("Error create hbox", err)
		return nil
	}
	box.SetVExpand(true)
	return box
}

func NewVBox(spacing int) *gtk.Box {
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, spacing)
	if err != nil {
		log.Println("Error create vbox", err)
		return nil
	}
	box.SetHExpand(true)
	return box
}
