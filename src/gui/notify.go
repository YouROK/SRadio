package gui

import (
	"github.com/gotk3/gotk3/glib"
)

func Notify(txt string) {
	nn := glib.NotificationNew("SRadio")
	nn.SetBody(txt)
	app.SendNotification("SRadio", nn)
}
