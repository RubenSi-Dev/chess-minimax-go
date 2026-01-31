package ai

import (
	"fmt"

	"github.com/spunker/chess/state"
)

// Heuristic evaluation functions, returning a score for a given state
// I could write more advanced heuristics, but for now I kept it simple

// getMaterialStat - (helperfunction) returns a map with the total material worth for each color
func getMaterialStat(s *state.State) (result map[string]int) {
	result = map[string]int{
		"white": 0,
		"black": 0,
	}

	for _, piece := range s.Board.GetPieces() {
		result[piece.Color] += piece.Worth
	}
	return
}

// evalMaterial - evaluates the material balance of the state (uses getMaterialStat)
// higher is better for white
func evalMaterial(s *state.State) int {
	stats := getMaterialStat(s)
	return stats["white"] - stats["black"]
}

// pieceMobility - (helperfunction) returns the mobility score of a piece (number of possible moves)
func pieceMobility(board *state.Board, piece *state.Piece) (int, error) {
	if piece == nil {
		return 0, fmt.Errorf("piece is nil")
	}
	if board == nil {
		return 0, fmt.Errorf("Board is nil")
	}
	possibleMoves, err := piece.GetPossibleMoves(board)
	if err != nil {
		return 0, err
	}
	return len(possibleMoves), nil

}

// evalMobility - evaluates the mobility balance of the state (higher is better for white)
func evalMobility(s *state.State) (int, error) {
	result := 0
	board := s.Board
	for _, piece := range board.GetPieces() {
		switch piece.Color {
		case "white":
			res, err := pieceMobility(s.Board, piece)
			if err != nil {
				return 0, err
			}
			result += res
		case "black":
			res, err := pieceMobility(s.Board, piece)
			if err != nil {
				return 0, err
			}
			result -= res
		}
	}
	return result, nil
}

// Weights - weights for different evaluation components
type Weights struct {
	Material float64
	Mobility float64
}

// EvalState - evaluates the state using weighted sum of different heuristics
func EvalState(s *state.State, weights *Weights) (float64, error) {
	evalMat := evalMaterial(s)
	evalMob, err := evalMobility(s)
	if err != nil {
		return 0, err
	}
	return weights.Material*float64(evalMat) + weights.Mobility*float64(evalMob), nil
}
