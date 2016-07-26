package gui

import (
	"github.com/gotk3/gotk3/glib"
)

var (
	app *glib.Application
)

func InitApp() {
	app = glib.ApplicationNew("ru.YouROK.SRadio", glib.APPLICATION_FLAGS_NONE)
}

func RunApp(f func()) {
	app.Connect("activate", f)
	app.Run(nil)
}
