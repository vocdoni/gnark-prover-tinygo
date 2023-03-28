# Gnark prover using TinyGo

> ⚠️🚧 This repository is currently a **proof of concept**. 🚧⚠️

This is an attempt to create an experimental zksnark prover/verifier for gnark circuits, compatible with [TinyGo](https://github.com/tinygo-org/tinygo), which means that it could be used on a bunch variety of targets (e.g. browsers).

It implements the same use case as [this circuit](https://github.com/vocdoni/zk-franchise-proof-circuit), and it attempts to replace the [Circom](https://github.com/iden3/circom) + [SnarkJS](https://github.com/iden3/snarkjs) stack.

## Motivations
 - 🚀 Gnark is very fast.
 - 💉 Supports unit testing and many modern backends and curves.
 - 🔗 [vocdoni-node](https://github.com/vocdoni/vocdoni-node) is currently writed enterily in Go, like Gnark, which will increase the maintainability of the source code.

## Project structure
```
  artifacts/        -> Includes generated artifacts for the implemented and compiled circuits.
  circuits/         -> Includes the available circuits definitions.
    zkcensus/       -> Port to gnark of https://github.com/vocdoni/zk-franchise-proof-circuit/blob/master/circuit/census.circom 
  cmd/compiler/     -> Simple command to compile available circuits.
  example/          -> Example of proof generation in js using gnark into a go-wasm.
  internal/
    circuit/        -> Internal definition of GenerateProof and VerifyProof funcs definitions.
      groth16/      -> Groth16 zk-snark backend version
      plonk/        -> Plonk zk-snark backend version
    zkaddress/      -> Alternative implementation of current vocdoni zkaddress (https://github.com/vocdoni/vocdoni-node/blob/master/crypto/zk/address.go)
  std/              -> Extended gnark std version with required ports.
    hash/poseidon/  -> Port to gnark of https://github.com/iden3/circomlib/blob/master/circuits/poseidon.circom
    smt/            -> Port to gnark of https://github.com/iden3/circomlib/blob/master/circuits/smt/smtverifier.circom
    zkaddress/      -> Implementation of the zkaddress using gnark
  wasm/             -> Wasm entrypoint and compiled version.
```

## Tests decription and results

### Followed stepts

1. Ports from Gnark of required Circom circuits:
    - `Poseidon hash`: [Gnark](./std/hash/poseidon/poseidon.go) | [Circom](https://github.com/iden3/circomlib/blob/db0202410771a3e3fc07c64c5226b64f954b8b5a/circuits/poseidon.circom).
    - `SMTVerifier`: [Gnark](./std/smt/verifier.go) | [Circom](https://github.com/iden3/circomlib/blob/a8cdb6cd1ad652cca1a409da053ec98f19de6c9d/circuits/smt/smtverifier.circom).
2. `ZkCensus` Vocdoni circuit port to Gnark: [Gnark](./circuits/zkcensus/zkcensus.go) | [Circom](https://github.com/vocdoni/zk-franchise-proof-circuit/blob/c2ead7f8502cf0dd7495140aec32599fd0a53199/circuit/census.circom).
4. `ZkCensus` (Gnark version) compiler and artifacts enconder command implementation.
    - Some blocks found and solved, read more [here](https://github.com/ConsenSys/gnark/issues/600).
3. Generic Gnark prover/verifier implementation.
    - Go WASM compiler as baseline. Found some incompatibilities with TinyGo, read more [here](https://github.com/tinygo-org/tinygo/issues/447#issuecomment-1455205919). 

### Requirements
* Go (1.20.2)
* TinyGo (@dgryski fork): dgryski/tinygo@a73e4c635331045f6d3cd49ddb0b9efd0019f94c


### Circuit 

#### Schema
```
                          +----+
  PUB_votingWeight+------>+ <= +------------------+--PRI_factoryWeight
                          +----+                  |
                                                  |
                          +-----------+           |
                          |           |           |
  PUB_censusRoot+-------->+           |(value)<---+
                          |           |
                          | SMT       |           +-----------+   +-----------+
                          | Verifier  |           |           |   |           |
  PRI_siblings+---------->+           |(key)<-----+ ZkAddress +<--+   pubKey  +---+-+PRI_privateKey
                          |           |           |           |   |           |   |
                          +-----------+           +-----------+   +-----------+   |
                                                                                  |
                                      +-----------+                               |
                          +----+      |           +<------------------------------+
  PUB_nullifier+--------->+ == +<-----+ Poseidon  |<------------+PUB_processID_0
                          +----+      |           +<------------+PUB_processID_1
                                      +-----------+
  PUB_voteHash
```

#### Inputs
| Name | Private/Public | Description |
|:---:|:---:|:---|
| *votingHeight* | `private` | The weight used to perform a vote. It must be equal to or lower than `factoryWeight`. |
| *factoryHeight* | `public` | The weight assigned to the voter as Merkle Tree leaf value. |
| *privateKey* | `private` | The voter private key. Seed of the ZkAddress.  |
| *censusRoot* | `public` | The Merkle Root of the current census tree. |
| *siblings* | `private` | Siblings of the voter ZkAddress leaf in the census tree. |
| *nullifier* | `public` | Parameter that combines the *privateKey* with the *electionId* to avoid proof reusability. |
| *electionId* | `public` | Encoded ID of the election. |
| *voteHash* | `public` | Parameter that combines the *privateKey* with the *factoryWeight* to be include it into the proof witness. |

### Available commands
* **Compile the prover and optimize the output**
  ```sh
  make compile-prover-{compiler}-{zk_backend}
  ```
  Select the desired WebAssembly compiler (`go` or `tinygo`) ZkSnark backend (`groth16` or `plonk`).

* **Compile the circuit artifacts**
  ```sh
  make compile-circuit-{zk_backend}
  ```
  Select the desired ZkSnark backend (`groth16` or `plonk`). It will override current artifacts.

* **Run example**
  ```sh
  make run-{compiler}-example
  ```
  Select the desired WebAssembly compiler (`go` or `tinygo`). It will use a previously compiled circuit artifacts.

### Code example

### Results

| Compiler | Snark Backend | Browser thread | Test result | Errors |
|:---:|:---:|:---:|:---:|:---:|
| Go (native) | Groth16 | main thread | ≈ 210s | ✅ |
| Go (native) | Plonk | main thread | ≈ 208s | ✅ |
| Go (native) | Groth16 | worker thread | ≈ 262s | ✅ |
| Go (native) | Plonk | worker thread | ≈ 211s | ✅ |
| TinyGo (dev) | Groth16 | main thread | - | ❌ `panic: reflect: unimplemented: AssignableTo with interface` |
| TinyGo (dev) | Plonk | main thread | - | ❌ `panic: reflect: unimplemented: AssignableTo with interface` |
| TinyGo (dev) | Groth16 | worker thread | - | ❌ `panic: reflect: unimplemented: AssignableTo with interface` |
| TinyGo (dev) | Plonk | worker thread | - | ❌ `panic: reflect: unimplemented: AssignableTo with interface` |
