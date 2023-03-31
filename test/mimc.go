package main

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

type TestMiMCCircuit struct {
	Input frontend.Variable `gnark:",public"`
}

func (c *TestMiMCCircuit) Define(api frontend.API) error {
	h, err := mimc.NewMiMC(api)
	if err != nil {
		panic(err)
	}

	var res frontend.Variable = 0
	for i := 0; i < 160; i++ {
		h.Write(res, c.Input)
		res = h.Sum()
		h.Reset()
	}
	return nil
}
