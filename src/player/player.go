package player

// #include <locale.h>
import "C"

import (
	"config"
	"log"
	"time"
	"unsafe"

	"github.com/yourok/go-mpv/mpv"
)

type EventCallback func(*mpv.Event, *Player)

type Player struct {
	mp        *mpv.Mpv
	isPlay    bool
	fileName  string
	cache     int
	cacheSeek int
	volume    int
	eventFunc EventCallback
}

func NewPlayer() *Player {
	p := &Player{}
	p.cache = config.GetCache()
	p.cacheSeek = config.GetCacheSeek()
	p.volume = 100
	return p
}

func (p *Player) Init() error {
	if p.mp == nil {
		buf := []byte("C")
		C.setlocale(C.LC_NUMERIC, (*C.char)(unsafe.Pointer(&buf[0])))
		m := mpv.Create()
		m.SetOption("no-resume-playback", mpv.FORMAT_FLAG, true)
		m.SetOption("volume", mpv.FORMAT_INT64, p.volume)
		m.SetOptionString("terminal", "yes")
		m.SetOptionString("softvol", "auto")
		m.SetOption("no-video", mpv.FORMAT_FLAG, true)
		m.SetOption("cache", mpv.FORMAT_INT64, p.cache)
		m.SetOption("cache-seek-min", mpv.FORMAT_INT64, p.cacheSeek)
		m.SetOption("volume", mpv.FORMAT_INT64, 0)
		err := m.Initialize()
		if err != nil {
			log.Println("Error init mpv", err)
			return err
		}
		p.mp = m
	}
	return nil
}

func (p *Player) Play(fileName string) error {
	if p.isPlay {
		return nil
	}
	var err error
	if fileName != p.fileName {
		p.Stop()
	}
	if p.mp == nil {
		err = p.Init()
	}

	if err != nil {
		return err
	}

	err = p.mp.Command([]string{"loadfile", fileName})
	if err != nil {
		p.Stop()
		return err
	}
	p.fileName = fileName
	p.isPlay = true
	go p.handler()
	return nil
}

func (p *Player) Stop() error {
	p.isPlay = false
	p.fileName = ""
	if p.mp != nil {
		p.mp.Command([]string{"stop"})
		p.mp.DetachDestroy()
		p.mp = nil
	}
	return nil
}

func (p *Player) Restart() error {
	fn := p.fileName
	p.Stop()
	return p.Play(fn)
}

func (p *Player) FadeIn() {
	//vol up
	for i := 0; i <= p.volume; i++ {
		if p.mp == nil {
			return
		}
		p.mp.SetProperty("volume", mpv.FORMAT_INT64, i)
		time.Sleep(time.Millisecond * 1)
	}
}

func (p *Player) FadeOut() {
	//vol down
	for i := p.volume; i >= 0; i-- {
		if p.mp == nil {
			return
		}
		p.mp.SetProperty("volume", mpv.FORMAT_INT64, i)
		time.Sleep(time.Millisecond * 1)
	}
}

func (p *Player) SetVolume(vol int) {
	if vol > 100 {
		vol = 100
	}
	if vol < 0 {
		vol = 0
	}
	p.volume = vol
	if p.mp == nil {
		return
	}
	p.mp.SetProperty("volume", mpv.FORMAT_INT64, vol)
}

func (p *Player) IsPlaying() bool {
	return p.isPlay
}

func (p *Player) SetEventCallback(fun EventCallback) {
	p.eventFunc = fun
}

func (p *Player) handler() {
	for p.isPlay {
		if p.eventFunc != nil {
			ev := p.mp.WaitEvent(-1)
			if ev.Event_Id != mpv.EVENT_NONE {
				p.eventFunc(ev, p)
			}
		}
		time.Sleep(time.Millisecond * 500)
	}
}
