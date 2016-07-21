package gui

import (
	"appindicator"
	"time"

	"github.com/gotk3/gotk3/gtk"
)

type TrayIcon struct {
	icon      *appindicator.AppIndicator
	iconsName []string
	menu      *gtk.Menu
	title     string
}

func NewTrayIcon() *TrayIcon {
	ti := &TrayIcon{}
	appind := appindicator.NewAppIndicator("example-simple-client", "indicator-messages", appindicator.CategoryOther)
	appind.SetStatus(appindicator.StatusActive)
	ti.icon = appind
	ti.title = "SRadio"
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
	t.icon.SetIcon(t.iconsName[index], t.title)
}

const (
	TA_PLAYING = iota
	TA_STOPPED
	TA_LOADING
	TA_CLOSE
)

type TrayAnimation struct {
	animIndex int
	frame     int
	trayicon  *TrayIcon
}

func NewTrayAnimation(ti *TrayIcon) *TrayAnimation {
	ta := &TrayAnimation{}
	ta.trayicon = ti
	ta.animIndex = TA_STOPPED
	go ta.animation()
	return ta
}

func (ta *TrayAnimation) Close() {
	ta.animIndex = TA_CLOSE
}

func (ta *TrayAnimation) SetAnimation(index int) {
	ta.animIndex = index
}

func (ta *TrayAnimation) GetTrayIcon() *TrayIcon {
	return ta.trayicon
}

func (ta *TrayAnimation) animation() {
	for ta.animIndex != TA_CLOSE {
		switch ta.animIndex {
		case TA_PLAYING:
			ta.animPlaying()
		case TA_LOADING:
			ta.animLoading()
		default:
			ta.animStoped()
		}
	}
}

func (ta *TrayAnimation) animPlaying() {
	ta.trayicon.SetIcon(ta.frame)
	ta.frame++
	if ta.frame == 3 {
		ta.frame = 0
		for i := 0; i < 30; i++ {
			time.Sleep(time.Second)
			if ta.animIndex != TA_PLAYING {
				break
			}
		}
	}
	time.Sleep(time.Second)
}

func (ta *TrayAnimation) animStoped() {
	ta.trayicon.SetIcon(ta.frame)
	ta.frame++
	if ta.frame == 2 {
		ta.frame = 0
	}
	time.Sleep(time.Second)
}

func (ta *TrayAnimation) animLoading() {
	ta.trayicon.SetIcon(ta.frame)
	ta.frame++
	if ta.frame == 3 {
		ta.frame = 0
		time.Sleep(time.Second)
	}
	time.Sleep(300 * time.Millisecond)
}
