package gui

import (
	"log"
	"os"
	"path/filepath"
	"radio"

	"github.com/gotk3/gotk3/gtk"
)

var (
	mainWnd     *gtk.Window
	radioEdit   []*gtk.Entry
	updwBtn     []*gtk.Button
	mainBtn     []*gtk.Button
	bottomBtn   []*gtk.Button
	radioList   *ListBox
	statusLabel *gtk.Label
	trayAnim    *TrayAnimation
	trayIcon    *TrayIcon

	rad *radio.Radio
)

func Init(r *radio.Radio) {
	var err error
	rad = r
	gtk.Init(nil)
	mainWnd, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Println("Error create window", err)
		return
	}

	mainWnd.SetPosition(gtk.WIN_POS_CENTER)
	mainWnd.SetTitle("SRadio")

	trayIcon = NewTrayIcon()
	if trayIcon != nil {
		icons := []string{filepath.Join(filepath.Dir(os.Args[0]), "radiotray1.png"),
			filepath.Join(filepath.Dir(os.Args[0]), "radiotray2.png"),
			filepath.Join(filepath.Dir(os.Args[0]), "radiotray3.png")}
		trayIcon.SetIconList(icons)
		trayIcon.SetIcon(0)

		trayAnim = NewTrayAnimation(trayIcon)
		trayAnim.SetAnimation(TA_STOPPED)
		mainWnd.Connect("destroy", func() {
			mainWnd.Hide()
		})
		updTrayMenu()

	} else {
		mainWnd.Connect("destroy", func() {
			gtk.MainQuit()
		})
	}

	mainWnd.Add(buildMainWnd())

	mainWnd.SetSizeRequest(600, 400)
	if trayIcon == nil {
		mainWnd.ShowAll()
	}

	go eventsHandler()

	gtk.Main()
	trayAnim.Close()
}

/*
********************************
* List R *      Add            *
* ------ *      Import         *
* ------ *      Export         *
* ------ *      Edit           *
* ------ *      Remove         *
* ------ *      Separator      *
****************************** *
*  Play  *  About  *  Exit     *
*  Stop  *         *           *
********************************
 */

func buildMainWnd() gtk.IWidget {
	mainVBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)
	{
		hbox := NewHBox(1)
		{ //Tree
			vboxleft := NewVBox(1)
			radioList = NewListBox()
			radioList.OnClick(radioListClick)
			radioList.OnDblClick(radioListDblClick)
			updateRadioList()

			hboxud := NewHBox(1)
			hboxud.SetVExpand(false)
			updwBtn := make([]*gtk.Button, 2)
			updwBtn[0], _ = gtk.ButtonNewWithLabel("Up")
			updwBtn[0].SetSizeRequest(-1, 25)
			updwBtn[1], _ = gtk.ButtonNewWithLabel("Down")
			updwBtn[1].SetSizeRequest(-1, 25)

			updwBtn[0].Connect("clicked", btnUp)
			updwBtn[1].Connect("clicked", btnDown)

			hboxud.PackStart(updwBtn[0], true, true, 1)
			hboxud.PackEnd(updwBtn[1], true, true, 1)

			scroll, _ := gtk.ScrolledWindowNew(nil, nil)
			scroll.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)
			scroll.Add(radioList.GetWidget())
			vboxleft.PackStart(scroll, true, true, 1)
			vboxleft.PackEnd(hboxud, false, false, 1)

			hbox.PackStart(vboxleft, true, true, 1)
		}
		{ // menu buttons
			radioEdit = make([]*gtk.Entry, 2)
			radioEdit[0], _ = gtk.EntryNew()
			radioEdit[1], _ = gtk.EntryNew()
			mainBtn = make([]*gtk.Button, 4)
			mainBtn[0], _ = gtk.ButtonNewWithLabel("Add")
			mainBtn[1], _ = gtk.ButtonNewWithLabel("Edit")
			mainBtn[2], _ = gtk.ButtonNewWithLabel("Remove")
			mainBtn[3], _ = gtk.ButtonNewWithLabel("Separator")
			vbox := NewVBox(2)
			for i := 0; i < len(radioEdit); i++ {
				vbox.PackStart(radioEdit[i], false, false, 1)
			}
			for i := 0; i < len(mainBtn); i++ {
				vbox.PackStart(mainBtn[i], true, true, 1)
			}
			hbox.Add(vbox)

			mainBtn[0].Connect("clicked", btnAdd)
			mainBtn[1].Connect("clicked", btnEdit)
			mainBtn[2].Connect("clicked", btnRemove)
			mainBtn[3].Connect("clicked", btnSeparator)

		}
		mainVBox.PackStart(hbox, true, true, 1)
	}
	{ //bottom buttons
		hbox := NewHBox(1)
		bottomBtn = make([]*gtk.Button, 3)
		bottomBtn[0], _ = gtk.ButtonNewWithLabel("Play")
		bottomBtn[1], _ = gtk.ButtonNewWithLabel("About")
		bottomBtn[2], _ = gtk.ButtonNewWithLabel("Exit")
		hbox.SetHomogeneous(true)
		for i := 0; i < len(bottomBtn); i++ {
			hbox.Add(bottomBtn[i])
			bottomBtn[i].SetSizeRequest(-1, 50)
		}
		mainVBox.PackStart(hbox, false, true, 15)

		bottomBtn[0].Connect("clicked", btnPlayClick)
		bottomBtn[1].Connect("clicked", btnAboutClick)
		bottomBtn[2].Connect("clicked", btnExitClick)
	}
	{ //status label
		statusLabel, _ = gtk.LabelNew("")
		statusLabel.SetJustify(gtk.JUSTIFY_LEFT)
		statusLabel.SetHAlign(gtk.ALIGN_START)
		mainVBox.PackEnd(statusLabel, false, true, 2)
		statusLabel.SetSizeRequest(-1, -1)
	}
	return mainVBox
}
