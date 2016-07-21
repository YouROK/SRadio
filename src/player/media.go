package player

import "github.com/yourok/go-mpv/mpv"

func (p *Player) GetMediaTitle() (string, error) {
	str, err := p.mp.GetProperty("media-title", mpv.FORMAT_STRING)
	return str.(string), err
}
