package state

func CreateMovePawn(white bool, from *Position, to *Position) *Move {
	if white && to.Y == 7 {
		return CreateMovePromotion(*from, *to, "queen")
	} else if !white && to.Y == 0 {
		return CreateMovePromotion(*from, *to, "queen")
	}
	return CreateMove(*from, *to)
}

func (p *Piece) getPossibleMovesPawn(board *Board) (possibleMoves []*Move) {
	possibleMoves = []*Move{}
	white := (p.Color == "white")

	var move0 Position
	if white {
		move0 = Position{X: p.Pos.X, Y: p.Pos.Y + 1}
	} else {
		move0 = Position{X: p.Pos.X, Y: p.Pos.Y - 1}
	}

	if board.GetPiece(&move0) == nil && board.isInBounds(&move0) {
		possibleMoves = append(possibleMoves, CreateMovePawn(white, &p.Pos, &move0))
		if (white && p.Pos.Y == 1) || (!white && p.Pos.Y == 6) {
			var move1 Position
			if white {
				move1 = Position{X: p.Pos.X, Y: p.Pos.Y + 2}
			} else {
				move1 = Position{X: p.Pos.X, Y: p.Pos.Y - 2}
			}
			if board.GetPiece(&move1) == nil {
				possibleMoves = append(possibleMoves, CreateMovePawn(white, &p.Pos, &move1))
			}
		}
	}

	var movesDiag []Position

	if white {
		movesDiag = []Position{
			{X: p.Pos.X + 1, Y: p.Pos.Y + 1},
			{X: p.Pos.X - 1, Y: p.Pos.Y + 1}}
	} else {
		movesDiag = []Position{
			{X: p.Pos.X + 1, Y: p.Pos.Y - 1},
			{X: p.Pos.X - 1, Y: p.Pos.Y - 1}}
	}

	for _, move := range movesDiag {
		if board.isInBounds(&move) {
			target := board.GetPiece(&move)
			if target != nil && target.Color != p.Color {
				possibleMoves = append(possibleMoves, CreateMovePawn(white, &p.Pos, &move))
			}
		}
	}
	return
}
