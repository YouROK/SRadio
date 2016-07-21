package radio

type RadioEvents struct {
	Type EventType
	Val  string
}

type EventType int

const (
	RE_START EventType = iota
	RE_STOP
	RE_ERROR
	RE_CHANGE_TITLE
	RE_PAUSE
	RE_UNPAUSE
)

func (e EventType) String() string {
	switch e {
	case RE_START:
		return "RE_START"
	case RE_STOP:
		return "RE_STOP"
	case RE_ERROR:
		return "RE_ERROR"
	case RE_CHANGE_TITLE:
		return "RE_CHANGE_TITLE"
	case RE_PAUSE:
		return "RE_PAUSE"
	case RE_UNPAUSE:
		return "RE_UNPAUSE"
	default:
		return "RE_UNKNOWN"
	}
}
