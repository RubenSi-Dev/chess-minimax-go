package main

import (
	"fmt"
	"slices"

	"github.com/spunker/chess/state"
)

type Game struct {
	State *state.State
	Over bool
	legalMovesPreProcess []*state.Move
}

func StartGame(setup string) (result *Game) {
	result = &Game{
		State: state.CreateState(setup),
		Over: false,
	}
	result.legalMovesPreProcess = result.State.GetLegalMoves()
	fmt.Println(result.legalMovesPreProcess)
	return
}

func (g *Game) PlayMoveAlgebraic(alg string) bool {
	return g.PlayMove(state.FromAlgebraicToMove(alg))
}

func (g *Game) PlayMove(move *state.Move) bool {
	legalMoves := g.State.GetLegalMoves()
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
