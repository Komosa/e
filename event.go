package main

import (
	"time"

	"github.com/nsf/termbox-go"
)

type ui struct {
	lines      [][]rune
	running    bool
	dontSave   bool
	sync, done chan struct{}
	ch         chan key
	f          map[key]func(*ui)
	cx, cy     int    // cursor pos
	prevkey    key    // for seqCont and vertMove
	prefpos    int    // for vertMove
	killbuf    []rune // for delete/un// for delete/undelete,
}

type key struct {
	mod termbox.Modifier
	key termbox.Key
	ch  rune
}

var hackKey = key{mod: termbox.ModAlt, ch: 'O'}

func newUI(data []byte) *ui {
	return &ui{
		lines:   fromBuf(data),
		running: true,
		sync:    make(chan struct{}),
		done:    make(chan struct{}),
		ch:      make(chan key, 16),
		f:       defMap(),
	}
}

// will be called in own goroutine
func (ui *ui) eventLoop() {
	for {
		e := termbox.PollEvent()
		switch e.Type {
		case termbox.EventResize:
			ui.sync <- struct{}{}
		case termbox.EventKey:
			if e.Key == termbox.KeyCtrlQ {
				close(ui.done)
				return
			}
			if e.Key == termbox.KeyCtrlRsqBracket {
				ui.dontSave = true
				close(ui.done)
				return
			}
			ui.ch <- key{key: e.Key, mod: e.Mod, ch: e.Ch}
		}
	}
}

func (ui *ui) dispatch(e key) {
	// fmt.Fprintf(os.Stderr, "dispatch: %#+v\n", e)
	if e.ch != 0 {
		e.key = 0
	}
	if f, ok := ui.f[e]; ok {
		f(ui)
	} else if e.ch != 0 && e.mod == 0 {
		ui.ins(e.ch)
	}
	ui.prevkey = e
}

func (ui *ui) update() {
	// until there is no event, there is no point to return from this func
	select {
	case <-ui.sync:
		return
	case <-ui.done:
		ui.running = false
		return
	case e := <-ui.ch:
		ui.dispatch(e)
	}

	// lets wait a moment before redraw
	for {
		select {
		case <-time.After(time.Millisecond * 15):
			return
		case <-ui.sync:
			return
		case <-ui.done:
			ui.running = false
			return
		case e := <-ui.ch:
			ui.dispatch(e)
		}
	}
}

func (ui *ui) trySeqCont(ch rune) bool {
	if ui.prevkey != hackKey {
		return false
	}
	switch ch {
	case 'c':
		fRightArrow(ui)
	case 'd':
		fLeftArrow(ui)
	case 'w':
		fHome(ui)
	case 'q':
		fEnd(ui)
	}
	return true
}
