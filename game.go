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
	result.legalMovesPreProcess, err = result.State.GetLegalMoves()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (g *Game) PlayMoveAlgebraic(alg string) (bool, error) {
	res, err := g.PlayMove(state.FromAlgebraicToMove(alg))
	return res, err
}

func (g *Game) PlayMove(move *state.Move) (bool, error) {
	legalMoves, err := g.State.GetLegalMoves()
	if err != nil {
		return false, err
	}
	var legalMove *state.Move
	isLegal := slices.ContainsFunc(legalMoves, func(m *state.Move) bool {
		if res, err := m.Equal(move); res {
			if err != nil {
				return false
			}
			legalMove = m
			return true
		}
		return false
	})

	if isLegal {
		_, err := g.State.ApplyMove(legalMove)
		if err != nil {
			return false, err
		}
		g.Moves++
		return true, nil
	}
	return false, nil
}

func (g *Game) String() string {
	return g.State.String()
}
