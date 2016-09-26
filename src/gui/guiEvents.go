package gui

import (
	"config"
	"log"
	"radio"

	"github.com/gotk3/gotk3/gtk"
)

func updateRadioList() {
	index := radioList.GetSelected()
	var items []string
	radios := config.GetRadios()
	for _, r := range radios {
		items = append(items, r.Name)
	}
	radioList.Update(items)
	radioList.SetSelected(index)
}

func eventsHandler() {
	radioChanEvent := rad.GetEvents()
	for true {
		radioEvent := <-radioChanEvent
		switch radioEvent.Type {
		case radio.RE_CHANGE_TITLE:
			{
				if radioEvent.Val != "" {
					mainWnd.SetTitle("SRadio - " + radioEvent.Val)
					statusLabel.SetLabel("Song: " + radioEvent.Val)
					trayIcon.SetTitle(radioEvent.Val)
					Notify(radioEvent.Val)
				}
			}
		case radio.RE_STOP:
			{
				mainWnd.SetTitle("SRadio")
				bottomBtn[0].SetLabel("Play")
				statusLabel.SetLabel("Stopped")
				if trayAnim != nil {
					trayAnim.SetAnimation(TA_STOPPED)
				}
			}
		case radio.RE_START, radio.RE_UNPAUSE:
			{
				if trayAnim != nil {
					trayAnim.SetAnimation(TA_PLAYING)
				}
				statusLabel.SetLabel("Playing... " + rad.GetRadio().Name)
			}
		case radio.RE_ERROR:
			{
				if trayAnim != nil {
					trayAnim.SetAnimation(TA_STOPPED)
				}
				log.Println("Error:", radioEvent.Val)
				bottomBtn[0].SetLabel("Play")
				statusLabel.SetLabel("Error: " + radioEvent.Val)
			}
		case radio.RE_PAUSE:
			{
				if trayAnim != nil {
					trayAnim.SetAnimation(TA_LOADING)
				}
				statusLabel.SetLabel("Loading...")
			}

		}
	}
}

func btnPlayClick() {
	if rad.IsPlaying() {
		rad.Stop()
	} else {
		rad.Play()
	}
	defer updTrayMenu()

	if mainWnd.IsVisible() {
		if rad.IsPlaying() {
			bottomBtn[0].SetLabel("Stop")
		} else {
			bottomBtn[0].SetLabel("Play")
		}
	}
}

func btnAboutClick() {
	about, _ := gtk.AboutDialogNew()
	about.SetVersion(config.Version)
	about.SetProgramName("SRadio")
	about.SetComments("SRadio - Simple Radio\nprogram for listen online radio\nAuthor: YouROK")
	about.SetWebsite("https://github.com/YouROK")
	about.Run()
	about.Hide()
}

func btnExitClick() {
	gtk.MainQuit()
}

func radioListClick(radInd int) {
	if radInd == -1 {
		return
	}
	radSel := config.GetRadios()[radInd]
	radioEdit[0].SetText(radSel.Name)
	radioEdit[1].SetText(radSel.Url)
}

func radioListDblClick(radInd int) {
	if radInd == -1 {
		return
	}
	defer updTrayMenu()
	radSel := config.GetRadios()[radInd]
	if radSel.Url == "separator" {
		return
	}
	config.SetSelectedRadio(radInd)
	if rad.GetRadio() == nil || radSel.Url != rad.GetRadio().Url {
		rad.Stop()
		rad.SetRadio(&radSel)
	}
	rad.Play()

	if mainWnd.IsVisible() {
		if rad.IsPlaying() {
			bottomBtn[0].SetLabel("Stop")
		} else {
			bottomBtn[0].SetLabel("Play")
		}
	}
}

func btnUp() {
	pos := radioList.GetSelected()
	if pos > 0 {
		up := config.GetRadios()[pos]
		upper := config.GetRadios()[pos-1]
		config.SetRadio(up, pos-1)
		config.SetRadio(upper, pos)
		updateRadioList()
		radioList.SetSelected(pos - 1)
	}
}

func btnDown() {
	pos := radioList.GetSelected()
	if pos < len(config.GetRadios())-1 {
		down := config.GetRadios()[pos]
		downer := config.GetRadios()[pos+1]
		config.SetRadio(down, pos+1)
		config.SetRadio(downer, pos)
		updateRadioList()
		radioList.SetSelected(pos + 1)
	}
}

func btnAdd() {
	rc := config.RadioCfg{}
	rc.Name, _ = radioEdit[0].GetText()
	rc.Url, _ = radioEdit[1].GetText()
	if rc.Name != "" && rc.Url != "" {
		config.AddRadio(rc)
		updateRadioList()
	}
}

func btnEdit() {
	rc := config.RadioCfg{}
	rc.Name, _ = radioEdit[0].GetText()
	rc.Url, _ = radioEdit[1].GetText()
	pos := radioList.GetSelected()
	if rc.Name != "" && rc.Url != "" && pos != -1 {
		config.SetRadio(rc, pos)
		updateRadioList()
	}
}

func btnRemove() {
	pos := radioList.GetSelected()
	if pos != -1 {
		config.DelRadio(pos)
		updateRadioList()
	}
}

func btnSeparator() {
	rc := config.RadioCfg{Name: "-------------------", Url: "separator"}
	config.AddRadio(rc)
	updateRadioList()
}

func updTrayMenu() {
	if trayIcon == nil {
		return
	}
	trayIcon.NewMenu()

	sel := config.GetSelectedRadio()
	defRadio := config.GetRadios()[sel]
	btnPlStName := "Play"
	if rad.IsPlaying() {
		btnPlStName = "Stop"
	}

	trayIcon.AddMenuItem(btnPlStName+" "+defRadio.Name, func(interface{}) { btnPlayClick() }, nil)
	trayIcon.AddMenuItem("Show settings", func(interface{}) {
		mainWnd.ShowAll()
		//TODO иногда не показывает контент
	}, nil)
	trayIcon.AddMenuItem("", nil, nil)

	for i, r := range config.GetRadios() {
		if r.Url == "separator" {
			trayIcon.AddMenuItem("", nil, nil)
		} else {
			trayIcon.AddMenuItem(r.Name, playRadio, i)
		}
	}
	trayIcon.AddMenuItem("", nil, nil)
	trayIcon.AddMenuItem("Exit", func(interface{}) { gtk.MainQuit() }, nil)
}

func playRadio(menuItm interface{}, i interface{}) {
	if ind, ok := i.(int); ok {
		radioListDblClick(ind)
	}
}
