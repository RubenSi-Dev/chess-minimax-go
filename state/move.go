package state

import (
	"fmt"
	"strings"
)

// Move - represents a move from one position to another
// has a promotion field for pawn promotions
type Move struct {
	From      Position
	To        Position
	Promotion string
}

// CreateMove - creates a new move from one position to another
func CreateMove(from Position, to Position) (result *Move) {
	result = &Move{
		From: from,
		To:   to,
	}
	return
}

// CreateMovePromotion - creates a new move with promotion
func CreateMovePromotion(from Position, to Position, prom string) (result *Move) {
	return &Move{
		From:      from,
		To:        to,
		Promotion: prom,
	}
}

func (this *Move) Equal(other *Move) bool {
	return this.From.Equal(other.From) && this.To.Equal(other.To)
}

func (m *Move) String() (result string) {
	result = fmt.Sprintf("{%v, %v", m.From, m.To)
	result += "}"
	return
}

// ToAlgebraic - converts the move to algebraic notation (e.g. e2-e4)
func (m *Move) ToAlgebraic() string {
	return fmt.Sprintf("%v-%v", m.From.ToAlgebraic(), m.To.ToAlgebraic())
}

// FromAlgebraicToMove - converts algebraic notation to a Move
func FromAlgebraicToMove(alg string) *Move {
	algs := strings.Split(alg, "-")
	return &Move{
		From: *fromAlgebraic(algs[0]),
		To:   *fromAlgebraic(algs[1]),
	}
}
