package ai

import (
	"math"

	"github.com/spunker/chess/state"
)

func minimax(s *state.State, depth int, max bool, alpha float64, beta float64, weights *Weights) (evaln float64) {
	if depth == 0 || s.IsGameOver() {
		return EvalState(s, weights)
	}


	if max { 
		evaln = math.Inf(-1) 
	} else {
		evaln = math.Inf(1) 
	}

	for _, move := range s.GetLegalMoves() {
		copyState := s.Copy()
		//fmt.Printf("bot evaluating move %v", move.ToAlgebraic())
		copyState.ApplyMove(move)
		currentEvaln := minimax(copyState, depth - 1, !max, alpha, beta, weights)
		if max {
			evaln = math.Max(evaln, currentEvaln)
			alpha = math.Max(alpha, evaln)
			if beta <= alpha { break }
		} else {
			evaln = math.Min(evaln, currentEvaln)
			beta = math.Min(beta, evaln)
			if beta <= alpha { break }
		}
	}
	return
}

func SelectMove(s *state.State, depth int, weights *Weights) (bestMove *state.Move, bestScore float64) {
	max := s.Turn == "white"
	if max {
		bestScore = math.Inf(-1)
	} else {
		bestScore = math.Inf(1)
	}
	for _, move := range s.GetLegalMoves() {
		copyState := s.Copy()
		copyState.ApplyMove(move)
		score := minimax(copyState, depth - 1, !max, math.Inf(-1), math.Inf(1), weights)
		if max {
			if score > bestScore {
				bestScore = score
				bestMove = move
			}
		} else {
			if score < bestScore {
				bestScore = score
				bestMove = move
			}
		}
	}
	return
}

