package state

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// Position - represents a position on the chess board
type Position struct {
	X int
	Y int
}

func (this Position) Equal(other Position) bool {
	return this.X == other.X &&
		this.Y == other.Y
}

func (p Position) String() (result string) {
	result = fmt.Sprintf("{%v, %v}", p.X, p.Y)
	return
}

func (p *Position) ToAlgebraic() (result string) {
	return fmt.Sprintf("%v%v", algebraicLetters[p.X], p.Y+1)
}

func fromAlgebraic(alg string) *Position {
	chars := strings.Split(alg, "")
	x := slices.Index(algebraicLetters, chars[0])
	y, _ := strconv.Atoi(chars[1])
	return &Position{
		X: x,
		Y: y - 1,
	}
}
