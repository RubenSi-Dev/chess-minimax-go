package state

var kingDirections = queenDirections

func (p Piece) getPossibleMovesKing(board *Board) ([]*Move, error) {
	result, err := p.stepMoves(board, &kingDirections)
	if err != nil {
		return nil, err
	}

	if res, err := p.isCastlingPossible(true, board); res {
		if err != nil {
			return nil, err
		}
		result = append(result, CreateMove(
			p.Pos,
			Position{X: p.Pos.X + 2, Y: p.Pos.Y},
		))
	}
	if res, err := p.isCastlingPossible(false, board); res {
		if err != nil {
			return nil, err
		}
		result = append(result, CreateMove(
			p.Pos,
			Position{X: p.Pos.X - 2, Y: p.Pos.Y},
		))
	}
	return result, nil
}

func (this Piece) isCastlingPossible(short bool, board *Board) (bool, error) {
	if this.HasMoved {
		return false, nil
	}
	posx := this.Pos.X
	posy := this.Pos.Y

	if short { // short castling
		rook, err := board.GetPiece(&Position{X: posx + 3, Y: posy})
		if err != nil {
			return false, err
		}

		p1, err1 := board.GetPiece(&Position{X: posx + 1, Y: posy})
		p2, err2 := board.GetPiece(&Position{X: posx + 2, Y: posy})
		if err1 != nil {
			return false, err1
		}
		if err2 != nil {
			return false, err2
		}

		if p1 != nil ||
			p2 != nil ||
			rook == nil ||
			rook.Type != "rook" ||
			rook.HasMoved {
			return false, nil
		}

		p1, err1 = board.GetPiece(&Position{X: posx + 1, Y: posy})
		p2, err2 = board.GetPiece(&Position{X: posx + 2, Y: posy})
		if err1 != nil {
			return false, err1
		}
		if err2 != nil {
			return false, err2
		}
		if p1 != nil ||
			p2 != nil ||
			rook == nil ||
			rook.Type != "rook" ||
			rook.HasMoved {
			return false, nil
		}
	} else { // long castling
		rook, err := board.GetPiece(&Position{X: posx - 4, Y: posy})
		if err != nil {
			return false, err
		}

		p1, err1 := board.GetPiece(&Position{X: posx - 1, Y: posy})
		p2, err2 := board.GetPiece(&Position{X: posx - 2, Y: posy})
		p3, err3 := board.GetPiece(&Position{X: posx - 3, Y: posy})
		if err1 != nil {
			return false, err1
		}
		if err2 != nil {
			return false, err2
		}
		if err3 != nil {
			return false, err2
		}
		if p1 != nil ||
			p2 != nil ||
			p3 != nil ||
			rook == nil ||
			rook.Type != "rook" ||
			rook.HasMoved {
			return false, nil
		}
	}

	return true, nil
}
