package smtverifier

import (
	"gnark-test/poseidon"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/bits"
)

func SMTVerifier(
	api frontend.API,
	enabled, root, oldKey, oldValue, isOld0, key, value, fnc frontend.Variable,
	siblings []frontend.Variable) error {
	nLevels := len(siblings)

	// Steps:
	//   1. Get the hash poseidon both, old and new key-value pair.
	//   2. Get the binary representation of both keys, old and new.
	//   3. Get the path of the current key.
	//   4. Calculate the root with the siblings provided.
	//   5. Compare the calculated root with the provided one.

	// [STEP 1]
	// hash1Old = H(oldKey | oldValue | 1)
	hash1Old := poseidon.Poseidon(api, []frontend.Variable{oldKey, oldValue, 1})
	// hash1New = H(key | value | 1)
	hash1New := poseidon.Poseidon(api, []frontend.Variable{key, value, 1})

	// [STEP 2]
	// component n2bOld = Num2Bits_strict();
	// n2bOld.in <== oldKey;
	// n2bOld := bits.ToBinary(api, oldKey)
	// component n2bNew = Num2Bits_strict();
	// n2bNew.in <== key;
	n2bNew := bits.ToBinary(api, key)

	// [STEP 3]
	// component smtLevIns = SMTLevIns(nLevels);
	// for (i=0; i<nLevels; i++) smtLevIns.siblings[i] <== siblings[i];
	// smtLevIns.enabled <== enabled;
	smtLevIns := SMTLevIns(api, siblings[:], enabled)
	// component sm[nLevels];
	sm := make([][5]frontend.Variable, nLevels)

	// sm[i] = SMTVerifierSM();
	st_top, st_iold, st_i0, st_inew, st_na := SMTVerifierSM(api,
		isOld0,              // sm[i].is0 <== isOld0
		smtLevIns[0],        // sm[i].levIns <== smtLevIns.levIns[i]
		fnc,                 // sm[i].fnc <== fnc
		enabled,             // sm[0].prev_top <== enabled
		0,                   // sm[0].prev_i0 <== 0
		0,                   // sm[0].prev_inew <== 0
		0,                   // sm[0].prev_iold <== 0
		api.Sub(1, enabled)) // sm[0].prev_na <== 1-enabled
	sm[0] = [5]frontend.Variable{st_top, st_iold, st_i0, st_inew, st_na}
	for i := 1; i < len(siblings); i++ {
		st_top, st_iold, st_i0, st_inew, st_na := SMTVerifierSM(api,
			isOld0,       // sm[i].is0 <== isOld0
			smtLevIns[i], // sm[i].levIns <== smtLevIns.levIns[i]
			fnc,          // sm[i].fnc <== fnc
			sm[i-1][0],   // sm[i].prev_top <== sm[i-1].st_top
			sm[i-1][2],   // sm[i].prev_i0 <== sm[i-1].st_i0
			sm[i-1][3],   // sm[i].prev_inew <== sm[i-1].st_inew
			sm[i-1][1],   // sm[i].prev_iold <== sm[i-1].st_iold
			sm[i-1][4])   // sm[i].prev_na <== sm[i-1].st_na
		sm[i] = [5]frontend.Variable{st_top, st_iold, st_i0, st_inew, st_na}
	}

	api.AssertIsEqual(api.Add(
		// sm[nLevels-1].st_na + sm[nLevels-1].st_iold + sm[nLevels-1].st_inew + sm[nLevels-1].st_i0 === 1
		sm[nLevels-1][4], sm[nLevels-1][1], sm[nLevels-1][3], sm[nLevels-1][2],
	), 1)

	// [STEP 4]
	levels := make([]frontend.Variable, nLevels)
	for i := nLevels - 1; i != -1; i-- {
		child := frontend.Variable(0)
		if i < nLevels-1 {
			child = levels[i+1]
		}

		levels[i] = SMTVerifierLevel(api,
			sm[i][0],    // st_top
			sm[i][2],    // st_i0
			sm[i][1],    // st_iold
			sm[i][3],    // st_new
			sm[i][4],    // st_na
			siblings[i], // sibling
			hash1Old,    // leaf1Old
			hash1New,    // leaf1New
			n2bNew[i],   // lrbit
			child,       // child
		)
	}

	// component areKeyEquals = IsEqual();
	// areKeyEquals.in[0] <== oldKey;
	// areKeyEquals.in[1] <== key;
	keysEqual := frontend.Variable(0)
	if api.Cmp(oldKey, key) == 0 {
		keysEqual = 1
	}

	// component keysOk = MultiAND(4);
	// keysOk.in[0] <== fnc;
	// keysOk.in[1] <== 1-isOld0;
	// keysOk.in[2] <== areKeyEquals.out;
	// keysOk.in[3] <== enabled;
	keysOk := multiAnd(api, fnc, api.Sub(1, isOld0), keysEqual, enabled)
	// keysOk.out === 0;
	api.AssertIsEqual(keysOk, 0)

	// [STEP 5]
	// component checkRoot = ForceEqualIfEnabled();
	// checkRoot.enabled <== enabled;
	// checkRoot.in[0] <== levels[0].root;
	// checkRoot.in[1] <== root;
	api.AssertIsEqual(root, levels[0])
	return nil
}

func SMTLevIns(api frontend.API, siblings []frontend.Variable, enabled frontend.Variable) []frontend.Variable {
	nLevels := len(siblings)
	// The last level must always have a sibling of 0. If not, then it cannot be inserted.
	// (isZero[nLevels-1].out - 1) * enabled === 0;
	if api.IsZero(enabled) == 0 {
		api.AssertIsEqual(siblings[nLevels-1], 0)
	}

	// for (i=0; i<nLevels; i++) {
	//     isZero[i] = IsZero();
	//     isZero[i].in <== siblings[i];
	// }
	isZero := make([]frontend.Variable, nLevels)
	isDone := make([]frontend.Variable, nLevels-1)
	for i := 0; i < nLevels; i++ {
		isZero[i] = api.IsZero(siblings[i])
	}

	levIns := make([]frontend.Variable, nLevels)
	last := api.Sub(1, isZero[nLevels-2])
	// levIns[nLevels-1] <== (1-isZero[nLevels-2].out);
	levIns[nLevels-1] = last
	// done[nLevels-2] <== levIns[nLevels-1];
	isDone[nLevels-2] = last
	for i := nLevels - 2; i > 0; i-- {
		// levIns[i] = (1-isDone[i])*(1-isZero[i-1])
		levIns[i] = api.Mul(api.Sub(1, isDone[i]), api.Sub(1, isZero[i-1]))
		// done[i-1] = levIns[i] + done[i]
		isDone[i-1] = api.Add(levIns[i], isDone[i])
	}
	// levIns[0] <== (1-done[0]);
	levIns[0] = api.Sub(1, isDone[0])

	return levIns
}

func SMTVerifierSM(api frontend.API,
	is0, levIns, fnc, prev_top, prev_i0, prev_iold, prev_inew, prev_na frontend.Variable) (
	st_top, st_iold, st_i0, st_inew, st_na frontend.Variable) {
	// prev_top_lev_ins <== prev_top * levIns;
	prevTopLevIns := api.Mul(prev_top, levIns)
	// prev_top_lev_ins_fnc <== prev_top_lev_ins*fnc
	prevTopLevInsFnc := api.Mul(prevTopLevIns, fnc)

	st_top = api.Sub(prev_top, prevTopLevIns)               // st_top <== prev_top - prev_top_lev_ins
	st_iold = api.Mul(prevTopLevInsFnc, api.Sub(1, is0))    // st_iold <== prev_top_lev_ins_fnc * (1 - is0)
	st_i0 = api.Mul(prevTopLevIns, is0)                     // st_i0 <== prev_top_lev_ins * is0;
	st_inew = api.Sub(prevTopLevIns, prevTopLevInsFnc)      // st_inew <== prev_top_lev_ins - prev_top_lev_ins_fnc
	st_na = api.Add(prev_na, prev_inew, prev_iold, prev_i0) // st_na <== prev_na + prev_inew + prev_iold + prev_i0
	return
}

// inputs:
//   - [0]: st_top;
//   - [1]: st_i0;
//   - [2]: st_iold;
//   - [3]: st_inew;
//   - [4]: st_na;
//   - [5]: sibling;
//   - [6]: old1leaf;
//   - [7]: new1leaf;
//   - [8]: lrbit;
//   - [9]: child;
func SMTVerifierLevel(api frontend.API,
	st_top, st_i0, st_iold, st_inew, st_na, sibling, old1leaf, new1leaf, lrbit,
	child frontend.Variable) frontend.Variable {
	// component switcher = Switcher();
	// switcher.sel <== lrbit;
	// switcher.L <== child;
	// switcher.R <== sibling;
	l, r := switcher(api, lrbit, child, sibling)
	// component proofHash = SMTHash2();
	// proofHash.L <== switcher.outL;
	// proofHash.R <== switcher.outR;
	proofHash := poseidon.Poseidon(api, []frontend.Variable{l, r})
	// aux[0] <== proofHash.out * st_top;
	aux0 := api.Mul(proofHash, st_top)
	// aux[1] <== old1leaf * st_iold;
	aux1 := api.Mul(old1leaf, st_iold)
	// root <== aux[0] + aux[1] + new1leaf * st_inew;
	return api.Add(aux0, aux1, api.Mul(new1leaf, st_inew))
}

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
