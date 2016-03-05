package main

import (
	mattnRuneWidth "github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

func runewidth(r rune) int {
	if r == '\t' {
		return 8
	}
	return mattnRuneWidth.RuneWidth(r)
}

func (ui *ui) draw() error {
	termbox.Clear(termbox.ColorDefault, termbox.ColorBlack)
	_, h := termbox.Size()
	for y, l := range ui.lines {
		if y == h {
			break
		}
		if len(l) > 0 && l[0] == '#' {
			ui.drawString(l, 0, y, termbox.ColorWhite, termbox.ColorBlack)
		} else {
			ui.drawString(l, 0, y, termbox.ColorDefault, termbox.ColorBlack)
		}

	}
	return termbox.Flush()
}

func (ui *ui) drawString(s []rune, x, y int, fg, bg termbox.Attribute) {
	w, _ := termbox.Size()
	for i, ch := range s {
		if y == ui.cy && i == ui.cx {
			termbox.SetCursor(x, ui.cy)
		}
		if ch == '\t' {
			x += 8
			continue
		}
		if ch == '\r' {
			if x+2 >= w {
				return
			}
			termbox.SetCell(x, y, '^', bg, fg)
			termbox.SetCell(x+1, y, 'R', bg, fg)
			x += 2
			continue
		}
		if x+runewidth(ch) >= w {
			return
		}
		termbox.SetCell(x, y, ch, fg, bg)
		x += runewidth(ch)
	}
	if y == ui.cy && len(s) == ui.cx {
		termbox.SetCursor(x, ui.cy)
	}
}
