package main

import "testing"

func TestCriarPlanilha(t *testing.T) {
	rq := Request{}
	rq.Aluno = "Guilherme Gomes Rocha"
	criarPlanilha(rq)
}
