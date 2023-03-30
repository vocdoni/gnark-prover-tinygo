package smt

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

// endLeafValue returns the encoded childless leaf value for the key-value pair
// provided, hashing it with the predefined hashing function 'H':
//
//	newLeafValue = H(key | value | 1)
func mimcEndLeafValue(api frontend.API, key, value frontend.Variable) frontend.Variable {
	h, err := mimc.NewMiMC(api)
	if err != nil {
		panic(err)
	}
	h.Write(key, value, 1)
	return h.Sum()
}

// func poseidonEndLeafValue(api frontend.API, key, value frontend.Variable) frontend.Variable {
// 	return poseidon.Hash(api, key, value, 1)
// }

// intermediateLeafValue returns the encoded intermediate leaf value for the
// key-value pair provided, hashing it with the predefined hashing function 'H':
//
//	intermediateLeafValue = H(l | r)
func mimcIntermediateLeafValue(api frontend.API, l, r frontend.Variable) frontend.Variable {
	h, err := mimc.NewMiMC(api)
	if err != nil {
		panic(err)
	}
	h.Write(l, r)
	return h.Sum()
}

// func poseidonIntermediateLeafValue(api frontend.API, l, r frontend.Variable) frontend.Variable {
// 	return poseidon.Hash(api, l, r)
// }

func switcher(api frontend.API, sel, l, r frontend.Variable) (outL, outR frontend.Variable) {
	// aux <== (R-L)*sel;
	aux := api.Mul(api.Sub(r, l), sel)
	// outL <==  aux + L;
	outL = api.Add(aux, l)
	// outR <== -aux + R;
	outR = api.Sub(r, aux)
	return
}

func multiAnd(api frontend.API, inputs ...frontend.Variable) frontend.Variable {
	if len(inputs) == 0 {
		return 0
	}
	if len(inputs) == 1 {
		return inputs[0]
	}

	res := inputs[0]
	for i := 1; i < len(inputs); i++ {
		res = api.And(res, inputs[i])
	}
	return res
}
