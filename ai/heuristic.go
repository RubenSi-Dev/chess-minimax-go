package ai

import (
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
func pieceMobility(board *state.Board, piece *state.Piece) int {
	return len(piece.GetPossibleMoves(board))

}

// evalMobility - evaluates the mobility balance of the state (higher is better for white)
func evalMobility(s *state.State) (result int) {
	board := s.Board
	for _, piece := range board.GetPieces() {
		switch piece.Color {
		case "white":
			result += pieceMobility(s.Board, piece)
		case "black":
			result -= pieceMobility(s.Board, piece)
		}
	}
	return
}

// Weights - weights for different evaluation components
type Weights struct {
	Material float64
	Mobility float64
}

// EvalState - evaluates the state using weighted sum of different heuristics
func EvalState(s *state.State, weights *Weights) (result float64) {
	return weights.Material*float64(evalMaterial(s)) + weights.Mobility*float64(evalMobility(s))
}
