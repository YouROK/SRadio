package gui

import (
	"appindicator"

	"github.com/gotk3/gotk3/gtk"
)

type TrayIcon struct {
	icon      *appindicator.AppIndicator
	iconsName []string
	menu      *gtk.Menu
	title     string
}

func NewTrayIcon(iconName string) *TrayIcon {
	ti := &TrayIcon{}
	appind := appindicator.NewAppIndicator("example-simple-client", "indicator-messages", appindicator.CategoryOther)
	appind.SetStatus(appindicator.StatusActive)
	appind.SetIcon(iconName, "SRadio")
	ti.icon = appind
	return ti
}

func (t *TrayIcon) NewMenu() {
	t.menu, _ = gtk.MenuNew()
	t.icon.SetMenu(t.menu)
}

func (t *TrayIcon) AddMenuItem(name string, callback interface{}, arg interface{}) {
	if name == "" {
		sep, _ := gtk.SeparatorMenuItemNew()
		sep.Show()
		t.menu.Append(sep)
	} else {
		itm1, _ := gtk.MenuItemNewWithLabel(name)
		itm1.Connect("activate", callback, arg)
		itm1.Show()
		t.menu.Append(itm1)
	}
}

func (t *TrayIcon) SetIconList(fileName []string) {
	t.iconsName = fileName
}

func (t *TrayIcon) SetTitle(title string) {
	t.title = title
	t.icon.SetTitle(title)
}

func (t *TrayIcon) SetIcon(index int) {
	t.icon.SetIcon(t.iconsName[index], t.iconsName[index])
}
