# Gnark prover using TinyGo

> ⚠️🚧 This repository is currently a **proof of concept**. 🚧⚠️

This is an experimental zksnark prover/verifier for gnark circuits that is compatible with TinyGo, which means that it could be used on a bunch variety of targets (e.g. browsers).

It implements the same use case as [this circuit](https://github.com/vocdoni/zk-franchise-proof-circuit), and it attempts to replace the Circom + SnarkJS stack.

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
    zkaddress/      -> Alternative implementation of current vocdoni zkaddress (https://github.com/vocdoni/vocdoni-node/blob/master/crypto/zk/address.go)
  std/              -> Extended gnark std version with required ports.
    hash/poseidon/  -> Port to gnark of https://github.com/iden3/circomlib/blob/master/circuits/poseidon.circom
    smt/            -> Port to gnark of https://github.com/iden3/circomlib/blob/master/circuits/smt/smtverifier.circom
    zkaddress/      -> Implementation of the zkaddress using gnark
  wasm/             -> Wasm entrypoint and compiled version.
```
