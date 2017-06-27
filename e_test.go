package e_test

import (
	"testing"

	e "github.com/Komosa/new_e"
)

func Test(t *testing.T) {
	f := func(in string, exp string) {
		out := e.Process(emit(in))
		if out != exp {
			t.Errorf("for %q:\n\tgot %q but\n\texpected %q", in, out, exp)
		}
	}

	f("", "")
	f("a", "a")            // entering text
	f("a←b", "ba")         // going back
	f("←a", "a")           // going on left border
	f("a←→b", "ab")        // going back and forth
	f("a→b", "ab")         // going on right border
	f("a\nb", "a\nb")      // entering new lines
	f("a\n↑b", "ba\n")     // going up (as in nano and geany)
	f("ab\nc↑d", "adb\nc") // entering in previous lines
	f("a↑b", "ba")         // going up in first line
	f("a\n↑b↓c", "ba\nc")  // going down
	f("a←↓b", "ab")        // going down in last line
}

var replace = map[rune]e.Event{
	'←':  {Key: e.KeyArrowLeft},
	'→':  {Key: e.KeyArrowRight},
	'↑':  {Key: e.KeyArrowUp},
	'↓':  {Key: e.KeyArrowDown},
	'\n': {Key: e.KeyEnter},
}

func emit(in string) []e.Event {
	var evs []e.Event
	for _, r := range in {
		if ev, ok := replace[r]; ok {
			evs = append(evs, ev)
		} else {
			evs = append(evs, e.Event{Ch: r})
		}
	}
	return evs
}