# Gnark prover using TinyGo

> ⚠️🚧 This repository is currently a **proof of concept**. 🚧⚠️

This is an attempt to create an experimental zksnark prover/verifier for gnark circuits, compatible with [TinyGo](https://github.com/tinygo-org/tinygo), which means that it could be used on a bunch variety of targets (e.g. browsers).

It implements the same use case as [this circuit](https://github.com/vocdoni/zk-franchise-proof-circuit), and it attempts to replace the [Circom](https://github.com/iden3/circom) + [SnarkJS](https://github.com/iden3/snarkjs) stack.

## Motivations
 - 🚀 [Gnark](https://github.com/ConsenSys/gnark) is very fast.
 - 💉 Supports unit testing and many modern backends and curves.
 - 🔗 [vocdoni-node](https://github.com/vocdoni/vocdoni-node) is currently writed enterily in Go, like Gnark, which will increase the maintainability of the source code.

**Read** about the comprehensive insights, decisions made, benchmarks and conclusions drawn in our article here https://hackmd.io/@vocdoni/B1VPA99Z3


## Project structure
```
  artifacts/        -> Includes generated artifacts for the implemented and compiled circuits.
  circuits/         -> Includes the available circuits definitions.
    zkcensus/       -> Port to gnark of https://github.com/vocdoni/zk-franchise-proof-circuit/blob/master/circuit/census.circom 
  cmd/
    compiler/       -> Simple command to compile available circuits.
    prover/         -> Simple command to test the prover with go.
  examples/         -> Example of proof generation in js using gnark into a go-wasm.
  internal/
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
    - Some blocks found and solved <sup>[1](#problems-found)</sup>.
3. Generic Gnark prover/verifier implementation.
    - Go WASM compiler as baseline. Found some incompatibilities with TinyGo<sup>[2](#problems-found)</sup>. 

### Requirements
* Go (1.20.2)
* TinyGo ([@vocdoni](https://github.com/vocdoni) fork): [vocdoni/tinygo](https://github.com/vocdoni/tinygo)


### Circuit 
The ZkCensus circuit anonymously proves that a voter is part of a census for a given election, without revealing its identity.

The ZkCensus circuit proves the following assertions:
1. The combination of the computed ZkAddress (using the given PrivateKey as seed) and the provided factoryWeight is a valid census tree leaf. This is tested computing the merkle root  with the candidate leaf and the provided siblings, and comparing the result with the provided census root.
2. The provided nullifier is valid. This is tested computing the nullifier with the electionID and the privateKey, and comparing the result with the provided nullifier.
3. The votingWeight is equal to or less than the factoryWeight.

Term descriptions:
* *ZkAddress*: the address of an anonymous voter in the Vochain, it is optimised for zk-snarks and helps to reduce the number of levels of the census merkle tree. It is based on the BN254 elliptic curve and it uses the voter private key as seed. Read more [here](https://github.com/vocdoni/vocdoni-node/blob/ca09fde59cef93f6b1de90c0c918adbff814e87e/crypto/zk/address.go).
* *Nullifier*: the result of applying the Poseidon hash to the combination of the election ID and the voter private key.

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
  PRI_siblings+---------->+           |(key)<-----+ ZkAddress +<--+ pubKey    +---+-+PRI_privateKey
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
| *votingHeight* | 🔐 `private` | The weight used to perform a vote. It must be equal to or lower than `factoryWeight`. |
| *factoryHeight* | 📢 `public` | The weight assigned to the voter as Merkle Tree leaf value. |
| *privateKey* | 🔐 `private` | The voter private key. Seed of the ZkAddress.  |
| *censusRoot* | 📢 `public` | The Merkle Root of the current census tree. |
| *siblings* | 🔐 `private` | Siblings of the voter ZkAddress leaf in the census tree. |
| *nullifier* | 📢 `public` | Parameter that combines the *privateKey* with the *electionId* to avoid proof reusability. |
| *electionId* | 📢 `public` | Encoded ID of the election. |
| *voteHash* | 📢 `public` | Parameter that combines the *privateKey* with the *factoryWeight* to be include it into the proof witness. |

### Available commands
* **Compile circuit, prover and run a web example**
  ```sh
  make run-tinygo-web-example-g16
  ```
  or for Plonk
  ```sh
  make run-tinygo-web-example-plonk
  ```

#### Other commands
See Makefile