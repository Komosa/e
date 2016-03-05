package main

func fromBuf(buf []byte) [][]rune {
	if len(buf) == 0 {
		return [][]rune{[]rune{}}
	}
	var lines [][]rune
	var line []rune
	nl := func() {
		lines = append(lines, append(line, []rune{}...))
		line = nil
	}
	for _, c := range string(buf) {
		if c == '\n' {
			nl()
		} else {
			line = append(line, c)
		}
	}
	if line != nil {
		nl()
	}
	return lines
}

func toBytes(lines [][]rune) []byte {
	for len(lines) > 0 {
		if len(lines[len(lines)-1]) == 0 {
			lines = lines[:len(lines)-1]
		} else {
			break
		}
	}
	var buf []byte
	for _, line := range lines {
		for _, c := range line {
			buf = append(buf, string(c)...)
		}
		buf = append(buf, '\n')
	}
	return buf
}
