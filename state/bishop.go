package state

var bishopDirections = []direction{
	{Dx: 1, Dy: 1},
	{Dx: 1, Dy: -1},
	{Dx: -1, Dy: -1},
	{Dx: -1, Dy: 1},
}

func (p *Piece) getPossibleMovesBishop(board *Board) ([]*Move, error) {
	raymoves, err := p.rayMoves(board, &bishopDirections)
	if err != nil {
		return nil, err
	}
	return raymoves, nil
}
