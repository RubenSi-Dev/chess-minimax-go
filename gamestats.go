package main

import (
	chess "github.com/spunker/chess/state"
)

type MaterialStat map[string]int

func GetMaterialStats(b *chess.Board) (result MaterialStat) {
	result = MaterialStat{
		"white": 0,
		"black": 0,
	}	
	
	for _, piece := range b.GetPieces() {
		result[piece.Color] += piece.Worth
	}
	return
}

func (m MaterialStat) GetAdvantage(color string) int {
	if color == "white" {
		return m["white"] - m["black"]
	} else {
		return m["black"] - m["white"]
	}
}
