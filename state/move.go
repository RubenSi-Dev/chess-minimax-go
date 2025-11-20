package state

import (
	"fmt"
	"strings"
) 

type Move struct {
	From Position
	To Position	
	//Promotion bool
	//Capture bool
}

func CreateMove(from Position, to Position) (result *Move) {
	result = &Move{
		From: from,
		To: to,
		//Promotion: prom,
		//Capture: capt,
	}
	return
}

func (this *Move) Equal(other *Move) bool {
	return this.From.Equal(other.From) && this.To.Equal(other.To)
}

func (m *Move) String() (result string) {
	result = fmt.Sprintf("{%v, %v", m.From, m.To)
	//if m.Promotion {
	//	result += ", Promotion"
	//}
	//if m.Capture {
		//result += ", Capture"
	//}
	result += "}"
	return
}

func (m *Move) ToAlgebraic() string {
	return fmt.Sprintf("%v-%v", m.From.ToAlgebraic(), m.To.ToAlgebraic())
}

func FromAlgebraicToMove(alg string) *Move {
	algs := strings.Split(alg, "-")	
	return &Move{
		From: *fromAlgebraic(algs[0]),
		To: *fromAlgebraic(algs[1]),
	}
}
