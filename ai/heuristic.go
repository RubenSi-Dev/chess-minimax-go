package ai

import (
	"github.com/spunker/chess/state"
)



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

func evalMaterial(s *state.State) int {
	stats := getMaterialStat(s)
	return stats["white"] - stats["black"]
}

func pieceMobility(board *state.Board, piece *state.Piece) int {
	return len(piece.GetPossibleMoves(board))

}

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

type Weights struct {
	Material float64 
	Mobility float64
}

func EvalState(s *state.State, weights *Weights) (result float64) {
	return weights.Material*float64(evalMaterial(s)) + weights.Mobility*float64(evalMobility(s))
}

