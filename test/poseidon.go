package main

import (
	"gnark-prover-tinygo/std/hash/poseidon"

	"github.com/consensys/gnark/frontend"
)

type TestPoseidonCircuit struct {
	Input frontend.Variable `gnark:",public"`
}

func (c *TestPoseidonCircuit) Define(api frontend.API) error {
	var res frontend.Variable = 0
	for i := 0; i < 160; i++ {
		res = poseidon.Hash(api, res, c.Input)
	}
	return nil
}
