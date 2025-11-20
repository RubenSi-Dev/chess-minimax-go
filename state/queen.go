package state

var queenDirections = append(rookDirections, bishopDirections...)

func (p Piece) getPossibleMovesQueen(board *Board) []*Move {
	return p.rayMoves(board, &queenDirections)
}
