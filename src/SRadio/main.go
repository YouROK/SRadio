package main

import (
	"config"
	"gui"
	"log"
	"radio"
	"runtime"
)

func main() {
	log.Println("Start SRadio", config.Version)
	runtime.LockOSThread()

	config.Init()
	radDef := config.GetSelectedRadio()
	rList := config.GetRadios()
	rad := radio.NewRadio()
	if radDef >= 0 && radDef < len(rList) {
		rad.SetRadio(&rList[radDef])
	} else if len(rList) > 0 {
		rad.SetRadio(&rList[0])
	}

	gui.Init(rad)
	rad.Stop()
}
