package main

import (
	"fmt"
	"github.com/spunker/chess/state"
)

func main() {
	//_test()
	//Interactive()
	StartTui()
}

func _test() {
	g := StartGame("default")
	
	fmt.Println(g)

	m1 := state.CreateMove(
		state.Position{X: 1, Y: 1},
		state.Position{X: 1, Y: 4},
	)
	fmt.Println(g.PlayMove(m1))
	fmt.Println(g)
	m1 = state.CreateMove(
		state.Position{X: 1, Y: 1},
		state.Position{X: 1, Y: 3},
	)
	fmt.Println(g.PlayMove(m1))
	fmt.Println(g)

	m1 = state.CreateMove(
		state.Position{X: 4, Y: 6},
		state.Position{X: 4, Y: 4},
	)
	fmt.Println(g.PlayMove(m1))
	fmt.Println(g)

	m1 = state.CreateMove(
		state.Position{X: 4, Y: 1},
		state.Position{X: 4, Y: 3},
	)
	fmt.Println(g.PlayMove(m1))
	fmt.Println(g)
}


