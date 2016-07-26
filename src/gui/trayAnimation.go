package gui

import (
	"time"
)

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
	isAnimate bool
	isReverse bool
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
	if ta.animIndex != index {
		ta.stopAnimation()
	}
	ta.animIndex = index
}

func (ta *TrayAnimation) GetTrayIcon() *TrayIcon {
	return ta.trayicon
}

func (ta *TrayAnimation) animation() {
	for ta.animIndex != TA_CLOSE {
		ta.isAnimate = true
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
	if ta.frame == 2 {
		ta.sleep(10000)
	}
	if ta.isReverse {
		ta.frame--
	} else {
		ta.frame++
	}

	if ta.frame < 1 || ta.frame > 1 {
		ta.isReverse = !ta.isReverse
	}

	ta.sleep(500)
}

func (ta *TrayAnimation) animStoped() {
	ta.trayicon.SetIcon(ta.frame)
	if ta.isReverse {
		ta.frame--
	} else {
		ta.frame++
	}
	if ta.frame < 1 || ta.frame >= 1 {
		ta.isReverse = !ta.isReverse
	}
	ta.sleep(1000)
}

func (ta *TrayAnimation) animLoading() {
	ta.trayicon.SetIcon(ta.frame)
	if ta.isReverse {
		ta.frame--
	} else {
		ta.frame++
	}
	if ta.frame < 1 || ta.frame > 1 {
		ta.isReverse = !ta.isReverse
	}
	ta.sleep(300)
}

func (t *TrayAnimation) sleep(millisecond int) {
	for i := 0; i < millisecond/100; i++ {
		if !t.isAnimate {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (t *TrayAnimation) stopAnimation() {
	t.isAnimate = false
}
