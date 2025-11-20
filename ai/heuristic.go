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

func EvalState(s *state.State) (result float64) {
	weights := map[string]float64 {
		"material": 2.0,
		"mobility": 0.5,
	}

	return weights["material"]*float64(evalMaterial(s)) + weights["mobility"]*float64(evalMobility(s))
}

