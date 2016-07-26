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
	updwBtn     []*Button
	mainBtn     []*Button
	bottomBtn   []*Button
	radioList   *ListBox
	statusLabel *gtk.Label
	trayAnim    *TrayAnimation
	trayIcon    *TrayIcon

	rad *radio.Radio
)

//TODO normal gui interface

func Init(r *radio.Radio) {
	rad = r
	gtk.Init(nil)
	InitApp()
	RunApp(initWnd)
}

func initWnd() {
	var err error
	mainWnd, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Println("Error create window", err)
		return
	}

	mainWnd.SetPosition(gtk.WIN_POS_CENTER)
	mainWnd.SetTitle("SRadio")

	trayIcon = NewTrayIcon(filepath.Join(filepath.Dir(os.Args[0]), "radiotray3.png"))
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

	app.Hold()
	gtk.Main()
	trayAnim.Close()
	app.Release()
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
			updwBtn := make([]*Button, 2)
			updwBtn[0] = NewButton("Up")
			updwBtn[0].SetSize(-1, 25)
			updwBtn[1] = NewButton("Down")
			updwBtn[1].SetSize(-1, 25)

			updwBtn[0].OnClick(btnUp)
			updwBtn[1].OnClick(btnDown)

			hboxud.PackStart(updwBtn[0].GetWidget(), true, true, 1)
			hboxud.PackEnd(updwBtn[1].GetWidget(), true, true, 1)

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
			mainBtn = make([]*Button, 4)
			mainBtn[0] = NewButton("Add")
			mainBtn[1] = NewButton("Edit")
			mainBtn[2] = NewButton("Remove")
			mainBtn[3] = NewButton("Separator")
			vbox := NewVBox(2)
			for i := 0; i < len(radioEdit); i++ {
				vbox.PackStart(radioEdit[i], false, false, 1)
			}
			for i := 0; i < len(mainBtn); i++ {
				vbox.PackStart(mainBtn[i].GetWidget(), true, true, 1)
			}
			hbox.Add(vbox)

			mainBtn[0].OnClick(btnAdd)
			mainBtn[1].OnClick(btnEdit)
			mainBtn[2].OnClick(btnRemove)
			mainBtn[3].OnClick(btnSeparator)

		}
		mainVBox.PackStart(hbox, true, true, 1)
	}
	{ //bottom buttons
		hbox := NewHBox(1)
		bottomBtn = make([]*Button, 3)
		bottomBtn[0] = NewButton("Play")
		bottomBtn[1] = NewButton("About")
		bottomBtn[2] = NewButton("Exit")
		hbox.SetHomogeneous(true)
		for i := 0; i < len(bottomBtn); i++ {
			hbox.Add(bottomBtn[i].GetWidget())
			bottomBtn[i].SetSize(-1, 50)
		}
		mainVBox.PackStart(hbox, false, true, 15)

		bottomBtn[0].OnClick(btnPlayClick)
		bottomBtn[1].OnClick(btnAboutClick)
		bottomBtn[2].OnClick(btnExitClick)
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
