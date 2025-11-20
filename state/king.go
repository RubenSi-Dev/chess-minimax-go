package state

var kingDirections = queenDirections

func (p Piece) getPossibleMovesKing(board *Board) []*Move {
	return p.stepMoves(board, &kingDirections)
}
