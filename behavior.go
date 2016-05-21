package main

import "github.com/nsf/termbox-go"

func (ui *ui) curLineLast() bool   { return ui.cy == len(ui.lines)-1 }
func (ui *ui) curRuneLast() bool   { return ui.cx == len(ui.lines[ui.cy]) }
func (ui *ui) posSpace(p int) bool { return isSpace(ui.lines[ui.cy][p]) }

func (ui *ui) ins(ch rune) {
	if !ui.trySeqCont(ch) {
		if isIndentEnd(ch) && ui.prev() == '\t' {
			fBackspace(ui)
		}
		ui.lines[ui.cy] = append(ui.lines[ui.cy][:ui.cx], append([]rune{ch}, ui.lines[ui.cy][ui.cx:]...)...)
		ui.cx++
	}
}

func fBackspace(ui *ui) {
	if ui.isCurBeg() {
		return
	}
	ui.horizMove(-1)
	fDelete(ui)
}

func fDelete(ui *ui) {
	if ui.curRuneLast() {
		if !ui.curLineLast() {
			ui.lines[ui.cy] = append(ui.lines[ui.cy], ui.lines[ui.cy+1]...)
			ui.lines = append(ui.lines[:ui.cy+1], ui.lines[ui.cy+2:]...)
		}
	} else {
		ui.lines[ui.cy] = append(ui.lines[ui.cy][:ui.cx], ui.lines[ui.cy][ui.cx+1:]...)
	}
}

func fEnd(ui *ui) { ui.cx = len(ui.lines[ui.cy]) }
func fHome(ui *ui) {
	px := ui.cx
	// smart home
	ui.cx = 0
	for _, c := range ui.lines[ui.cy] {
		if c == ' ' || c == '\t' {
			ui.cx++
		} else {
			break
		}
	}
	if ui.cx == px {
		// smarter home
		ui.cx = 0
	}
}

func fEnter(ui *ui) {
	nl := []rune{}
	for i, ch := range ui.lines[ui.cy] {
		if i >= ui.cx || ch != '\t' {
			break
		}
		nl = append(nl, ch)
	}
	if isIndentBeg(ui.prev()) {
		nl = append(nl, '\t')
	}
	nx := len(nl)
	if !ui.curRuneLast() {
		nl = append(nl, ui.lines[ui.cy][ui.cx:]...)
		ui.lines[ui.cy] = ui.lines[ui.cy][:ui.cx]
	}
	ui.lines = append(ui.lines[:ui.cy+1], append([][]rune{nl}, ui.lines[ui.cy+1:]...)...)
	ui.cy++
	ui.cx = nx
}

func delWord(ui *ui) {
	n := len(ui.lines[ui.cy])
	if n == 0 {
		return
	}
	if ui.cx == n {
		ui.cx--
	}
	d := 0
	// delete also left part of word if invoked from middle
	if !ui.posSpace(ui.cx) {
		for ui.cx > 0 && !ui.posSpace(ui.cx-1) {
			ui.cx--
			d++
		}
	} else {
		// delete spaces from cursor to next word
		for ui.cx+d+1 < n && ui.posSpace(ui.cx+d+1) {
			d++
		}
	}
	// delete from cursor to end of word
	for ui.cx+d+1 < n && !ui.posSpace(ui.cx+d+1) {
		d++
	}
	// delete also all but one space after word
	after := 0
	for ui.cx+d+1 < n && ui.posSpace(ui.cx+d+1) {
		d++
		after = 1
	}
	d -= after
	ui.lines[ui.cy] = append(ui.lines[ui.cy][:ui.cx], ui.lines[ui.cy][ui.cx+d+1:]...)
}

func delLine(ui *ui) {
	ui.killbuf = ui.lines[ui.cy]
	ui.cx = 0
	if len(ui.lines) == 1 {
		ui.lines = [][]rune{[]rune{}}
	} else if ui.curLineLast() {
		ui.lines[ui.cy] = []rune{}
	} else {
		ui.lines = append(ui.lines[:ui.cy], ui.lines[ui.cy+1:]...)
	}
}
func undelLine(ui *ui) {
	if ui.killbuf == nil {
		return
	}
	b := make([]rune, len(ui.killbuf))
	copy(b, ui.killbuf)
	ui.lines = append(ui.lines[:ui.cy], append([][]rune{b}, ui.lines[ui.cy:]...)...)
	ui.cy++
	ui.cx = 0
}

func (ui *ui) vertMove(dir int) {
	switch ui.prevkey {
	case key{key: termbox.KeyArrowUp}, key{key: termbox.KeyArrowDown}:
		// calculate preffered pos on move initiation
	default:
		ui.prefpos = 0
		for i, r := range ui.lines[ui.cy] {
			if i == ui.cx {
				break
			}
			ui.prefpos += runewidth(r)
		}
	}
	if ui.cy+dir >= 0 && ui.cy+dir < len(ui.lines) {
		ui.cy += dir
	}
	pos := 0
	ui.cx = 0
	for _, r := range ui.lines[ui.cy] {
		if pos >= ui.prefpos {
			break
		}
		pos += runewidth(r)
		ui.cx++
	}
}

func (ui *ui) horizMove(dir int) {
	if ui.cx+dir < 0 || ui.cx+dir > len(ui.lines[ui.cy]) {
		if ui.cy+dir >= 0 && ui.cy+dir < len(ui.lines) {
			ui.cy += dir
			if dir == 1 {
				ui.cx = 0
			} else {
				ui.cx = len(ui.lines[ui.cy])
			}
		}
	} else {
		ui.cx += dir
	}
}

// character under cursor
func (ui *ui) cur() rune {
	if ui.curRuneLast() {
		return '\n'
	}
	return ui.lines[ui.cy][ui.cx]
}

// character before cursor
func (ui *ui) prev() rune {
	if ui.cx == 0 {
		return '\n'
	}
	return ui.lines[ui.cy][ui.cx-1]
}
func isSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}
func isIndentBeg(ch rune) bool {
	return ch == '{' || ch == '(' || ch == '[' || ch == ':'
}
func isIndentEnd(ch rune) bool {
	return ch == '}' || ch == ')' || ch == ']'
}
func (ui *ui) isCurEnd() bool {
	return ui.curRuneLast() && ui.curLineLast()
}
func (ui *ui) isCurBeg() bool {
	return ui.cx == 0 && ui.cy == 0
}

func fRightArrow(ui *ui) {
	for !ui.isCurEnd() && isSpace(ui.cur()) {
		ui.horizMove(1)
	}
	for !ui.isCurEnd() && !isSpace(ui.cur()) {
		ui.horizMove(1)
	}
}
func fLeftArrow(ui *ui) {
	for !ui.isCurBeg() && isSpace(ui.cur()) {
		ui.horizMove(-1)
	}
	for !ui.isCurBeg() && !isSpace(ui.cur()) {
		ui.horizMove(-1)
	}
}

func defMap() map[key]func(*ui) {
	return map[key]func(*ui){
		key{ch: 'd', mod: termbox.ModAlt}: delWord,
		key{key: termbox.KeyDelete}:       fDelete,
		key{key: termbox.KeyBackspace}:    fBackspace,
		key{key: termbox.KeyBackspace2}:   fBackspace,
		key{key: termbox.KeyTab}:          func(ui *ui) { ui.ins('\t') },
		key{key: termbox.KeyEnter}:        fEnter,
		key{key: termbox.KeySpace}:        func(ui *ui) { ui.ins(' ') },
		key{key: termbox.KeyHome}:         fHome,
		key{key: termbox.KeyEnd}:          fEnd,
		key{key: termbox.KeyArrowUp}:      func(ui *ui) { ui.vertMove(-1) },
		key{key: termbox.KeyArrowDown}:    func(ui *ui) { ui.vertMove(1) },
		key{key: termbox.KeyArrowLeft}:    func(ui *ui) { ui.horizMove(-1) },
		key{key: termbox.KeyArrowRight}:   func(ui *ui) { ui.horizMove(1) },
		key{key: termbox.KeyCtrlK}:        delLine,
		key{key: termbox.KeyCtrlU}:        undelLine,
		key{key: termbox.KeyCtrlO}: func(ui *ui) {
			delLine(ui)
			fEnter(ui)
			ui.horizMove(-1)
		},
		key{key: termbox.KeyCtrlE}: func(ui *ui) {
			ui.lines = [][]rune{[]rune{}}
			ui.cx = 0
			ui.cy = 0
		},
	}
}
