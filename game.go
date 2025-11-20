package main

import (
	"github.com/spunker/chess/state"
	//"fmt"
	"slices"
)

type Game struct {
	State *state.State
	Over bool
}

func StartGame(setup string) *Game {
	return &Game{
		State: state.CreateState(setup),
		Over: false,
	}
}

func (g *Game) PlayMoveAlgebraic(alg string) bool {
	return g.PlayMove(state.FromAlgebraicToMove(alg))
}


func (g *Game) PlayMove(move *state.Move) bool {
	//fmt.Println(move.ToAlgebraic())
	legalMoves := g.State.GetLegalMoves()
	//fmt.Println(legalMoves)
	var legalMove *state.Move
	isLegal := slices.ContainsFunc(legalMoves, func(m *state.Move) bool {
		if m.Equal(move) {
			legalMove = m
			return true
		}
		return false
	}) 

	if isLegal {
		g.State.ApplyMove(legalMove)
		return true
	}
	return false
}

func (g *Game) String() string {
	return g.State.String()
}
