package main

import (
	"bytes"
	"unicode/utf8"

	"github.com/zyedidia/tcell"
)

func runeToByteIndex(n int, txt []byte) int {
	if n == 0 {
		return 0
	}

	count := 0
	i := 0
	for len(txt) > 0 {
		_, size := utf8.DecodeRune(txt)

		txt = txt[size:]
		count += size
		i++

		if i == n {
			break
		}
	}
	return count
}

type Line struct {
	text   []byte
	colors []tcell.Style
}

// A LineArray simply stores and array of lines and makes it easy to insert
// and delete in it
type LineArray struct {
	lines []Line
}

func (la *LineArray) setLines(text [][]byte) {
	for i, line := range la.lines {
		line.text = text[i]
	}
}

// NewLineArray returns a new line array from an array of bytes
func NewLineArray(text []byte) *LineArray {
	la := new(LineArray)
	// Split the bytes into lines
	split := bytes.Split(text, []byte("\n"))
	la.lines = make([]Line, len(split))
	// la.setLines(make([][]byte, len(split)))
	for i := range split {
		la.lines[i].text = make([]byte, len(split[i]))
		copy(la.lines[i].text, split[i])
	}

	return la
}

// Returns the String representation of the LineArray
func (la *LineArray) String() string {
	str := ""
	for _, line := range la.lines {
		str += string(line.text) + "\n"
	}
	return str
	// return string(bytes.Join(la.lines, []byte("\n")))
}

// NewlineBelow adds a newline below the given line number
func (la *LineArray) NewlineBelow(y int) {
	la.lines = append(la.lines, Line{[]byte(" "), []tcell.Style{}})
	copy(la.lines[y+2:], la.lines[y+1:])
	la.lines[y+1] = Line{[]byte(""), []tcell.Style{}}
}

// inserts a byte array at a given location
func (la *LineArray) insert(pos Loc, value []byte) {
	x, y := runeToByteIndex(pos.X, la.lines[pos.Y].text), pos.Y
	// x, y := pos.x, pos.y
	for i := 0; i < len(value); i++ {
		if value[i] == '\n' {
			la.Split(Loc{x, y})
			x = 0
			y++
			continue
		}
		la.insertByte(Loc{x, y}, value[i])
		x++
	}
}

// inserts a byte at a given location
func (la *LineArray) insertByte(pos Loc, value byte) {
	la.lines[pos.Y].text = append(la.lines[pos.Y].text, 0)
	copy(la.lines[pos.Y].text[pos.X+1:], la.lines[pos.Y].text[pos.X:])
	la.lines[pos.Y].text[pos.X] = value
}

// JoinLines joins the two lines a and b
func (la *LineArray) JoinLines(a, b int) {
	la.insert(Loc{len(la.lines[a].text), a}, la.lines[b].text)
	la.DeleteLine(b)
}

// Split splits a line at a given position
func (la *LineArray) Split(pos Loc) {
	la.NewlineBelow(pos.Y)
	la.insert(Loc{0, pos.Y + 1}, la.lines[pos.Y].text[pos.X:])
	la.DeleteToEnd(Loc{pos.X, pos.Y})
}

// removes from start to end
func (la *LineArray) remove(start, end Loc) string {
	sub := la.Substr(start, end)
	startX := runeToByteIndex(start.X, la.lines[start.Y].text)
	endX := runeToByteIndex(end.X, la.lines[end.Y].text)
	if start.Y == end.Y {
		la.lines[start.Y].text = append(la.lines[start.Y].text[:startX], la.lines[start.Y].text[endX:]...)
	} else {
		for i := start.Y + 1; i <= end.Y-1; i++ {
			la.DeleteLine(start.Y + 1)
		}
		la.DeleteToEnd(Loc{startX, start.Y})
		la.DeleteFromStart(Loc{endX - 1, start.Y + 1})
		la.JoinLines(start.Y, start.Y+1)
	}
	return sub
}

// DeleteToEnd deletes from the end of a line to the position
func (la *LineArray) DeleteToEnd(pos Loc) {
	la.lines[pos.Y].text = la.lines[pos.Y].text[:pos.X]
}

// DeleteFromStart deletes from the start of a line to the position
func (la *LineArray) DeleteFromStart(pos Loc) {
	la.lines[pos.Y].text = la.lines[pos.Y].text[pos.X+1:]
}

// DeleteLine deletes the line number
func (la *LineArray) DeleteLine(y int) {
	la.lines = la.lines[:y+copy(la.lines[y:], la.lines[y+1:])]
}

// DeleteByte deletes the byte at a position
func (la *LineArray) DeleteByte(pos Loc) {
	la.lines[pos.Y].text = la.lines[pos.Y].text[:pos.X+copy(la.lines[pos.Y].text[pos.X:], la.lines[pos.Y].text[pos.X+1:])]
}

// Substr returns the string representation between two locations
func (la *LineArray) Substr(start, end Loc) string {
	startX := runeToByteIndex(start.X, la.lines[start.Y].text)
	endX := runeToByteIndex(end.X, la.lines[end.Y].text)
	if start.Y == end.Y {
		return string(la.lines[start.Y].text[startX:endX])
	}
	var str string
	str += string(la.lines[start.Y].text[startX:]) + "\n"
	for i := start.Y + 1; i <= end.Y-1; i++ {
		str += string(la.lines[i].text) + "\n"
	}
	str += string(la.lines[end.Y].text[:endX])
	return str
}
