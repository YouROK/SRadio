package radio

import (
	"config"
	"log"
	"player"

	"github.com/yourok/go-mpv/mpv"
)

type Callback func(ev *mpv.Event, val interface{})

type Radio struct {
	pl        *player.Player
	playing   *config.RadioCfg
	events    chan RadioEvents
	songTitle string
	errors    int
}

func NewRadio() *Radio {
	r := &Radio{}
	r.pl = player.NewPlayer()
	r.events = make(chan RadioEvents, 255)
	r.pl.SetEventCallback(r.eventsHandler)
	return r
}

func (r *Radio) SetRadio(rcfg *config.RadioCfg) {
	r.playing = rcfg
}

func (r *Radio) GetRadio() *config.RadioCfg {
	return r.playing
}

func (r *Radio) GetEvents() chan RadioEvents {
	return r.events
}

func (r *Radio) Play() {
	if r.playing != nil {
		r.songTitle = ""
		r.errors = 0
		r.pl.Play(r.playing.Url)
		r.sendEvent(RE_START, "")
	}
}

func (r *Radio) Stop() {
	r.pl.FadeOut()
	r.pl.Stop()
	r.sendEvent(RE_STOP, "")
}

func (r *Radio) Restart() {
	r.pl.FadeOut()
	r.pl.Stop()
	r.songTitle = ""
	r.pl.Play(r.playing.Url)
}

func (r *Radio) IsPlaying() bool {
	return r.pl.IsPlaying()
}

func (r *Radio) eventsHandler(ev *mpv.Event, pl *player.Player) {
	if ev.Event_Id == mpv.EVENT_START_FILE {
		pl.FadeIn()
	}
	if ev.Event_Id == mpv.EVENT_END_FILE {
		if r.errors > 5 {
			r.pl.Stop()
			r.errors = 0
			log.Println("So many errors, stop and send error")
			var val string
			if data, ok := ev.Data.(mpv.EventEndFile); ok {
				val = data.ErrCode.Error()
			}
			r.sendEvent(RE_ERROR, val)
		}

		if r.IsPlaying() {
			r.errors++
			r.Restart()
			log.Println("Restart player", r.errors)
			return
		}
	}

	if ev.Event_Id == mpv.EVENT_METADATA_UPDATE {
		title, err := r.pl.GetMediaTitle()
		if err == nil && r.songTitle != title {
			r.sendEvent(RE_CHANGE_TITLE, title)
			r.songTitle = title
		}
	}

	if ev.Event_Id == mpv.EVENT_PAUSE {
		r.sendEvent(RE_PAUSE, "")
	}
	if ev.Event_Id == mpv.EVENT_UNPAUSE {
		r.sendEvent(RE_UNPAUSE, "")
	}
}

func (r *Radio) sendEvent(t EventType, val string) {
	r.events <- RadioEvents{t, val}
}
