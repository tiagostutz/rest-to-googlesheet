package main

import "testing"

func TestCreateSheet(t *testing.T) {
	rq := Request{}
	rq.Title = "John Doe"
	criarPlanilha(rq)
}
