package ai

import (
	"fmt"
	"math"
	"sync"

	"github.com/spunker/chess/state"
)

type Result struct {
	move  *state.Move
	score float64
	err   error
}

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
		return math.Inf(-1), err
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

	results := make(chan Result, len(legalMoves))
	var wg sync.WaitGroup

	for _, move := range legalMoves {
		wg.Add(1)
		go evaluateMove(s, move, depth, weights, max, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		if res.err != nil {
			return nil, math.Inf(-1), res.err
		}
		if max {
			if res.score > bestScore {
				bestScore = res.score
				bestMove = res.move
			}
		} else {
			if res.score < bestScore {
				bestScore = res.score
				bestMove = res.move
			}
		}
	}
	return bestMove, bestScore, nil
}

func evaluateMove(s *state.State, m *state.Move, depth int, weights *Weights, max bool, wg *sync.WaitGroup, results chan<- Result) {
	defer wg.Done()
	copyState, err := s.Copy()
	if err != nil {
		results <- Result{move: m, score: math.Inf(-1), err: err}
	}

	_, err = copyState.ApplyMove(m)
	if err != nil {
		results <- Result{move: m, score: math.Inf(-1), err: err}
	}

	score, err := minimax(copyState, depth-1, !max, math.Inf(-1), math.Inf(1), weights)
	if err != nil {
		results <- Result{move: m, score: score, err: fmt.Errorf("error evalutating move %v", m.ToAlgebraic())}
	}

	results <- Result{move: m, score: score, err: nil}
}
