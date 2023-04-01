# Test MiMC vs Poseidon over Gnark

> ⚠️🚧 This repository is currently a **proof of concept**. 🚧⚠️


This branch implement (check the [`/test`](./test)) a simple benchmark to test the performance of MiMC hash versus Poseidon into a [Gnark](https://github.com/ConsenSys/gnark) circuit.

It uses the MiMC implementation of ConsenSys ([source code](https://github.com/ConsenSys/gnark/blob/master/std/hash/mimc/mimc.go)) and our Poseidon implementation ([source code](./std/hash/poseidon/poseidon.go)).


### Requirements
* Go (1.20.2)

### Available commands
* **Compile and run the test**
  ```sh
  make make run-mimc-poseidon-test
  ```

### Results

| Platform | Backend | MiMC | Poseidon |
|:---|:---:|---:|---:|
| Go | Groth16 | 1.260s | 0.327s |
| Go | Plonk | 0.295s | 0.227s |
| Browser | Groth16 | 235.82s | 86.56s |
| Browser | Plonk | 268.57s | 139.87s |

```
Macmini9,1 (Z12N0004MY/A), Chip Apple M1 (8 cores), 16 GB Memory
Google Chrome Versión 111.0.5563.146 (Build oficial) (arm64)
```