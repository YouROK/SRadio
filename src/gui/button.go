package gui

import (
	"log"
	"sync"

	"github.com/gotk3/gotk3/gtk"
)

type Button struct {
	btn   *gtk.Button
	mutex sync.Mutex
}

func NewButton(label string) *Button {
	btn := &Button{}
	var err error
	btn.btn, err = gtk.ButtonNewWithLabel(label)
	if err != nil {
		log.Println("Error create button", err)
		return nil
	}

	return btn
}

func (b *Button) SetLabel(label string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.btn.SetLabel(label)
}

func (b *Button) SetSize(width, height int) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.btn.SetSizeRequest(width, height)
}

func (b *Button) OnClick(f func()) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.btn.Connect("clicked", f)
}

func (b *Button) GetWidget() *gtk.Button {
	return b.btn
}
