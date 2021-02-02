package main

const (
	clearEOL = "\033[K"
	moveUp   = "\033[A"
	margin   = 4
)

var console struct {
	w int
	h int
}
