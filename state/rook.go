package state

var rookDirections = []direction{
	{Dx: 1, Dy: 0},
	{Dx: -1, Dy: 0},
	{Dx: 0, Dy: 1},
	{Dx: 0, Dy: -1},
}

func (p Piece) getPossibleMovesRook(board *Board) []*Move {
	return p.rayMoves(board, &rookDirections)
}
