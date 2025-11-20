package state

var bishopDirections = []direction{
	{Dx: 1, Dy: 1},
	{Dx: 1, Dy: -1},
	{Dx: -1, Dy: -1},
	{Dx: -1, Dy: 1},
}

func (p *Piece) getPossibleMovesBishop(board *Board) []*Move {
	return p.rayMoves(board, &bishopDirections)
}
