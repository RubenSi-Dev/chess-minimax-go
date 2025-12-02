package ai

import (
	"fmt"
	"math"

	"github.com/spunker/chess/state"
)


// minimax - minimax algorithm with alpha-beta pruning
func minimax(s *state.State, depth int, max bool, alpha float64, beta float64, weights *Weights) (evaln float64, err error) {
	err = nil

	// base case for recursion
	if depth == 0 || s.IsGameOver() {
		return EvalState(s, weights), err
	}


	// check if maximizing or minimizing player	
	if max { 
		evaln = math.Inf(-1) 
	} else {
		evaln = math.Inf(1) 
	}

	// iterate through all legal moves fetched from the state
	for _, move := range s.GetLegalMoves() {
		// try the move (simulate on a copy)
		copyState := s.Copy()
		copyState.ApplyMove(move)

		// recursively call minimax on the new state
		currentEvaln, err := minimax(copyState, depth - 1, !max, alpha, beta, weights)
		if err != nil {
			return currentEvaln, fmt.Errorf("error evalutating %v", move.ToAlgebraic())
		}

		// update evaln, alpha, beta based on maximizing or minimizing player
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

// SelectMove - selects the best move using minimax algorithm
func SelectMove(s *state.State, depth int, weights *Weights) (bestMove *state.Move, bestScore float64, err error) {
	err = nil

	// hardcoded for maximizing player being white (for now)
	max := s.Turn == "white"
	if max {
		bestScore = math.Inf(-1)
	} else {
		bestScore = math.Inf(1)
	}

	// iterate through all legal moves this basically does the first layer of minimax because minimax itself doesn't return the move
	for _, move := range s.GetLegalMoves() {
		copyState := s.Copy()
		copyState.ApplyMove(move)
		score, err := minimax(copyState, depth - 1, !max, math.Inf(-1), math.Inf(1), weights)
		if err != nil {
			return nil, score, fmt.Errorf("error evalutating move %v", move.ToAlgebraic())
		}
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
	fmt.Printf("selected move %v", bestMove.ToAlgebraic())
	return 
}

