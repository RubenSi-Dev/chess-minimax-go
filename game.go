package main

import (
	"slices"

	"github.com/spunker/chess/state"
)

type Game struct {
	State                *state.State
	Over                 bool
	Moves                int
	legalMovesPreProcess []*state.Move
}

func StartGame(setup string) (*Game, error) {
	newState, err := state.CreateState(setup)
	if err != nil {
		return nil, err
	}
	result := &Game{
		State: newState,
		Over:  false,
		Moves: 0,
	}
	result.legalMovesPreProcess = result.State.GetLegalMoves()
	return result, nil
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
		g.Moves++
		return true
	}
	return false
}

func (g *Game) String() string {
	return g.State.String()
}
