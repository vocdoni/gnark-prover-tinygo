// smt is a port of Circom SMTVerifier. It attempts to check a proof
// of a Sparse Merkle Tree (compatible with Arbo Merkle Tree implementation).
// Check the original implementation from Iden3:
//   - https://github.com/iden3/circomlib/tree/a8cdb6cd1ad652cca1a409da053ec98f19de6c9d/circuits/smt
package smt

import "github.com/consensys/gnark/frontend"

func Verifier(api frontend.API,
	root, key, value frontend.Variable, siblings []frontend.Variable) error {
	return smtverifier(api, 1, root, 0, 0, 0, key, value, 0, siblings)
}

func smtverifier(api frontend.API,
	enabled, root, oldKey, oldValue, isOld0, key, value, fnc frontend.Variable,
	siblings []frontend.Variable) error {
	nLevels := len(siblings)

	// Steps:
	//   1. Get the hash of both key-value pairs, old and new one.
	//   2. Get the binary representation of the key new.
	//   3. Get the path of the current key.
	//   4. Calculate the root with the siblings provided.
	//   5. Compare the calculated root with the provided one.

	// [STEP 1]
	// hash1Old = H(oldKey | oldValue | 1)
	hash1Old := mimcEndLeafValue(api, oldKey, oldValue)
	// hash1New = H(key | value | 1)
	hash1New := mimcEndLeafValue(api, key, value)

	// [STEP 2]
	// component n2bNew = Num2Bits_strict();
	// n2bNew.in <== key;
	n2bNew := api.ToBinary(key, api.Compiler().FieldBitLen())

	// [STEP 3]
	// component smtLevIns = SMTLevIns(nLevels);
	// for (i=0; i<nLevels; i++) smtLevIns.siblings[i] <== siblings[i];
	// smtLevIns.enabled <== enabled;
	levIns := smtLevIns(api, siblings[:], enabled)
	// component sm[nLevels];
	sm := make([][5]frontend.Variable, nLevels)

	// sm[i] = SMTVerifierSM();
	stTop, stIold, stI0, stInew, stNa := smtVerifierSM(api,
		isOld0,              // sm[i].is0 <== isOld0
		levIns[0],           // sm[i].levIns <== smtLevIns.levIns[i]
		fnc,                 // sm[i].fnc <== fnc
		enabled,             // sm[0].prev_top <== enabled
		0,                   // sm[0].prev_i0 <== 0
		0,                   // sm[0].prev_inew <== 0
		0,                   // sm[0].prev_iold <== 0
		api.Sub(1, enabled)) // sm[0].prev_na <== 1-enabled
	sm[0] = [5]frontend.Variable{stTop, stIold, stI0, stInew, stNa}
	for i := 1; i < len(siblings); i++ {
		stTop, stIold, stI0, stInew, stNa := smtVerifierSM(api,
			isOld0,     // sm[i].is0 <== isOld0
			levIns[i],  // sm[i].levIns <== smtLevIns.levIns[i]
			fnc,        // sm[i].fnc <== fnc
			sm[i-1][0], // sm[i].prev_top <== sm[i-1].st_top
			sm[i-1][2], // sm[i].prev_i0 <== sm[i-1].st_i0
			sm[i-1][3], // sm[i].prev_inew <== sm[i-1].st_inew
			sm[i-1][1], // sm[i].prev_iold <== sm[i-1].st_iold
			sm[i-1][4]) // sm[i].prev_na <== sm[i-1].st_na
		sm[i] = [5]frontend.Variable{stTop, stIold, stI0, stInew, stNa}
	}

	// sm[nLevels-1].st_na + sm[nLevels-1].st_iold + sm[nLevels-1].st_inew + sm[nLevels-1].st_i0 === 1
	api.AssertIsEqual(api.Add(
		sm[nLevels-1][4], sm[nLevels-1][1], sm[nLevels-1][3], sm[nLevels-1][2],
	), 1)

	// [STEP 4]
	levels := make([]frontend.Variable, nLevels)
	for i := nLevels - 1; i != -1; i-- {
		child := frontend.Variable(0)
		if i < nLevels-1 {
			child = levels[i+1]
		}

		levels[i] = smtVerifierLevel(api,
			sm[i][0],    // st_top
			sm[i][1],    // st_iold
			sm[i][3],    // st_new
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

func smtLevIns(api frontend.API, siblings []frontend.Variable, enabled frontend.Variable) []frontend.Variable {
	nLevels := len(siblings)
	// The last level must always have a sibling of 0. If not, then it cannot be inserted.
	// (isZero[nLevels-1].out - 1) * enabled === 0;
	if api.IsZero(enabled) == 0 {
		api.AssertIsEqual(siblings[nLevels-1], 0)
	}
	isDone := make([]frontend.Variable, nLevels-1)
	levIns := make([]frontend.Variable, nLevels)
	last := api.Sub(1, api.IsZero(siblings[nLevels-2]))
	// levIns[nLevels-1] <== (1-isZero[nLevels-2].out);
	levIns[nLevels-1] = last
	// done[nLevels-2] <== levIns[nLevels-1];
	isDone[nLevels-2] = last
	for i := nLevels - 2; i > 0; i-- {
		// isZero[i-1] = IsZero();
		// isZero[i-1].in <== siblings[i];
		// levIns[i] = (1-isDone[i])*(1-isZero[i-1])
		levIns[i] = api.Mul(api.Sub(1, isDone[i]), api.Sub(1, api.IsZero(siblings[i-1])))
		// done[i-1] = levIns[i] + done[i]
		isDone[i-1] = api.Add(levIns[i], isDone[i])
	}
	// levIns[0] <== (1-done[0]);
	levIns[0] = api.Sub(1, isDone[0])
	return levIns
}

func smtVerifierSM(api frontend.API,
	is0, levIns, fnc, prevTop, prevI0, prevIold, prevInew, prevNa frontend.Variable) (
	stTop, stIold, stI0, stInew, stNa frontend.Variable) {
	// prev_top_lev_ins <== prev_top * levIns;
	prevTopLevIns := api.Mul(prevTop, levIns)
	// prev_top_lev_ins_fnc <== prev_top_lev_ins*fnc
	prevTopLevInsFnc := api.Mul(prevTopLevIns, fnc)

	stTop = api.Sub(prevTop, prevTopLevIns)             // st_top <== prev_top - prev_top_lev_ins
	stIold = api.Mul(prevTopLevInsFnc, api.Sub(1, is0)) // st_iold <== prev_top_lev_ins_fnc * (1 - is0)
	stI0 = api.Mul(prevTopLevIns, is0)                  // st_i0 <== prev_top_lev_ins * is0;
	stInew = api.Sub(prevTopLevIns, prevTopLevInsFnc)   // st_inew <== prev_top_lev_ins - prev_top_lev_ins_fnc
	stNa = api.Add(prevNa, prevInew, prevIold, prevI0)  // st_na <== prev_na + prev_inew + prev_iold + prev_i0
	return
}

func smtVerifierLevel(api frontend.API, stTop, stIold, stInew, sibling,
	old1leaf, new1leaf, lrbit, child frontend.Variable) frontend.Variable {
	// component switcher = Switcher();
	// switcher.sel <== lrbit;
	// switcher.L <== child;
	// switcher.R <== sibling;
	l, r := switcher(api, lrbit, child, sibling)
	// component proofHash = SMTHash2();
	// proofHash.L <== switcher.outL;
	// proofHash.R <== switcher.outR;
	proofHash := mimcIntermediateLeafValue(api, l, r)
	// aux[0] <== proofHash.out * st_top;
	aux0 := api.Mul(proofHash, stTop)
	// aux[1] <== old1leaf * st_iold;
	aux1 := api.Mul(old1leaf, stIold)
	// root <== aux[0] + aux[1] + new1leaf * st_inew;
	return api.Add(aux0, aux1, api.Mul(new1leaf, stInew))
}
