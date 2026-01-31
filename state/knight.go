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

func (p Piece) getPossibleMovesKnight(board *Board) ([]*Move, error) {
	raymoves, err := p.stepMoves(board, &knightDirections)
	if err != nil {
		return nil, err
	}
	return raymoves, nil
}
