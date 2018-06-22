package e

import "strings"

type Event struct {
	Key uint16
	Ch  rune
}

func Process(evs []Event) string {
	var txt string
	var cur int

	var prev_was_updown bool    // was prev key up or down arrow?
	var pref_col, pref_line int // preffered cur pos after sequence of up and down arrows

	// insert next key as txt[cur]; cur++
	put := func(ch rune) {
		txt = txt[:cur] + string(ch) + txt[cur:]
		cur++
	}

	for _, ev := range evs {
		if ev.Key == KeyArrowUp || ev.Key == KeyArrowDown {
			if !prev_was_updown {
				pref_line = strings.Count(txt[:cur], "\n")
				line_beg := 1 + strings.LastIndex(txt[:cur], "\n")
				pref_col = cur - line_beg
			}

			if ev.Key == KeyArrowUp {
				pref_line--
			} else {
				pref_line++
			}
			prev_was_updown = true
		} else {
			if prev_was_updown {
				prev_was_updown = false
				cur = 0

				// go to pref line
				for p := 0; p < pref_line; p++ {
					jump := strings.Index(txt[cur:], "\n") + 1
					if jump == 0 {
						cur = len(txt)
						break
					}
					cur += jump
				}

				// go to pref col
				if pref_line >= 0 {
					for p := 0; cur < len(txt) && p < pref_col && txt[cur] != '\n'; p++ {
						cur++
					}
				}
			}

			switch {
			case ev.Ch != 0:
				put(ev.Ch)
			case ev.Key == KeyEnter:
				put('\n')
			case ev.Key == KeyArrowLeft:
				cur = max(cur-1, 0)
			case ev.Key == KeyArrowRight:
				cur = min(cur+1, len(txt))
			}
		}
	}

	return txt
}

func min(a, b int) int {
	if a > b {
		a = b
	}
	return a
}
func max(a, b int) int {
	if a < b {
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
