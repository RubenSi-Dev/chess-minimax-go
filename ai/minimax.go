package ai

import (
	"fmt"
	"math"

	"github.com/spunker/chess/state"
)

// minimax - minimax algorithm with alpha-beta pruning
func minimax(s *state.State, depth int, max bool, alpha float64, beta float64, weights *Weights) (float64, error) {
	// base case for recursion
	isOver, err := s.IsGameOver()
	if err != nil {
		return 0, err
	}
	if depth == 0 || isOver {
		eval, err := EvalState(s, weights)
		if err != nil {
			return 0, err
		}
		return eval, nil
	}

	// check if maximizing or minimizing player
	var evaln float64
	if max {
		evaln = math.Inf(-1)
	} else {
		evaln = math.Inf(1)
	}

	// iterate through all legal moves fetched from the state
	legalMoves, err := s.GetLegalMoves()
	if err != nil {
		return 0, err
	}
	for _, move := range legalMoves {
		// try the move (simulate on a copy)
		copyState, err := s.Copy()
		if err != nil {
			return evaln, err
		}
		_, err = copyState.ApplyMove(move)
		if err != nil {
			return evaln, err
		}

		// recursively call minimax on the new state
		currentEvaln, err := minimax(copyState, depth-1, !max, alpha, beta, weights)
		if err != nil {
			return currentEvaln, fmt.Errorf("error evalutating %v", move.ToAlgebraic())
		}

		// update evaln, alpha, beta based on maximizing or minimizing player
		if max {
			evaln = math.Max(evaln, currentEvaln)
			alpha = math.Max(alpha, evaln)
			if beta <= alpha {
				break
			}
		} else {
			evaln = math.Min(evaln, currentEvaln)
			beta = math.Min(beta, evaln)
			if beta <= alpha {
				break
			}
		}
	}
	return evaln, nil
}

// SelectMove - selects the best move using minimax algorithm
func SelectMove(s *state.State, depth int, weights *Weights) (*state.Move, float64, error) {
	// hardcoded for maximizing player being white (for now)
	var bestMove *state.Move
	var bestScore float64
	max := s.Turn == "white"
	if max {
		bestScore = math.Inf(-1)
	} else {
		bestScore = math.Inf(1)
	}

	// iterate through all legal moves this basically does the first layer of minimax because minimax itself doesn't return the move
	legalMoves, err := s.GetLegalMoves()
	if err != nil {
		return nil, bestScore, err
	}
	for _, move := range legalMoves {
		copyState, err := s.Copy()
		if err != nil {
			return nil, bestScore, err
		}

		_, err = copyState.ApplyMove(move)
		if err != nil {
			return nil, bestScore, err
		}

		score, err := minimax(copyState, depth-1, !max, math.Inf(-1), math.Inf(1), weights)
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
	//fmt.Printf("selected move %v", bestMove.ToAlgebraic())
	return bestMove, bestScore, nil
}
