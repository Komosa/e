package e

import "strings"

type Event struct {
	Key uint16
	Ch  rune
}

func Process(evs []Event) string {
	var out string
	var cur int

	put := func(ch rune) {
		out = out[:cur] + string(ch) + out[cur:]
		cur++
	}
	currLineNum := func() int {
		return strings.Count(out[:cur], "\n")
	}
	startOfLine := func(nth int) int {
		if nth <= 0 {
			return -1
		}
		for i, c := range out {
			if c == '\n' {
				nth--
				if nth == 0 {
					return i
				}
			}
		}
		return len(out) + nth
	}

	goLineUp := func() {
		cl := currLineNum()
		if cl == 0 {
			cur = 0
			return
		}
		sl := startOfLine(cl)
		cur += startOfLine(cl-1) - sl
		cur = min(cur, sl)
	}
	goLineDown := func() {
		cl := currLineNum()
		cur += startOfLine(cl+1) - startOfLine(cl)
		cur = min(cur, len(out))
	}

	for _, ev := range evs {
		switch {
		case ev.Ch != 0:
			put(ev.Ch)
		case ev.Key == KeyEnter:
			put('\n')
		case ev.Key == KeyArrowLeft:
			if cur != 0 {
				cur--
			}
		case ev.Key == KeyArrowRight:
			cur = min(cur+1, len(out))
		case ev.Key == KeyArrowUp:
			goLineUp()
		case ev.Key == KeyArrowDown:
			goLineDown()
		}
	}

	return out
}

func min(a, b int) int {
	if a > b {
		a = b
	}
	return a
}

const (
	KeyF1 uint16 = 0xFFFF - iota
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyInsert
	KeyDelete
	KeyHome
	KeyEnd
	KeyPgup
	KeyPgdn
	KeyArrowUp
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight

	MouseLeft
	MouseMiddle
	MouseRight
	MouseRelease
	MouseWheelUp
	MouseWheelDown

	KeyCtrlTilde      = 0x00
	KeyCtrl2          = 0x00
	KeyCtrlSpace      = 0x00
	KeyCtrlA          = 0x01
	KeyCtrlB          = 0x02
	KeyCtrlC          = 0x03
	KeyCtrlD          = 0x04
	KeyCtrlE          = 0x05
	KeyCtrlF          = 0x06
	KeyCtrlG          = 0x07
	KeyBackspace      = 0x08
	KeyCtrlH          = 0x08
	KeyTab            = 0x09
	KeyCtrlI          = 0x09
	KeyCtrlJ          = 0x0A
	KeyCtrlK          = 0x0B
	KeyCtrlL          = 0x0C
	KeyEnter          = 0x0D
	KeyCtrlM          = 0x0D
	KeyCtrlN          = 0x0E
	KeyCtrlO          = 0x0F
	KeyCtrlP          = 0x10
	KeyCtrlQ          = 0x11
	KeyCtrlR          = 0x12
	KeyCtrlS          = 0x13
	KeyCtrlT          = 0x14
	KeyCtrlU          = 0x15
	KeyCtrlV          = 0x16
	KeyCtrlW          = 0x17
	KeyCtrlX          = 0x18
	KeyCtrlY          = 0x19
	KeyCtrlZ          = 0x1A
	KeyEsc            = 0x1B
	KeyCtrlLsqBracket = 0x1B
	KeyCtrl3          = 0x1B
	KeyCtrl4          = 0x1C
	KeyCtrlBackslash  = 0x1C
	KeyCtrl5          = 0x1D
	KeyCtrlRsqBracket = 0x1D
	KeyCtrl6          = 0x1E
	KeyCtrl7          = 0x1F
	KeyCtrlSlash      = 0x1F
	KeyCtrlUnderscore = 0x1F
	KeySpace          = 0x20
	KeyBackspace2     = 0x7F
	KeyCtrl8          = 0x7F
)
