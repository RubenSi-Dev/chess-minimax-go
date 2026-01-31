package state

var rookDirections = []direction{
	{Dx: 1, Dy: 0},
	{Dx: -1, Dy: 0},
	{Dx: 0, Dy: 1},
	{Dx: 0, Dy: -1},
}

func (p Piece) getPossibleMovesRook(board *Board) ([]*Move, error) {
	raymoves, err := p.rayMoves(board, &rookDirections)
	if err != nil {
		return nil, err
	}
	return raymoves, nil
}
