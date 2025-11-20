package state

var knightDirections = []direction{
	{Dx: 2, Dy: 1},
	{Dx: -2, Dy: 1},
	{Dx: 2, Dy: -1},
	{Dx: -2, Dy: -1},
	{Dx: 1, Dy: 2},
	{Dx: -1, Dy: 2},
	{Dx: 1, Dy: -2},
	{Dx: -1, Dy: -2},
}

func (p Piece) getPossibleMovesKnight(board *Board) []*Move {
	return p.stepMoves(board, &knightDirections)
}
