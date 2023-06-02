import './wasm_exec.js';

import { init, WASI } from '@wasmer/wasi';

import witnessUrl from "url:./artifacts/zkcensus.witness";
import wasmUrl from "url:./artifacts/g16_prover.wasm";

addEventListener('message', async(e) => {
    await init();

    let start = Date.now();
    console.log("reading artifacts...");

    let resWitness = await fetch(witnessUrl);
    let witnessBuff = await resWitness.arrayBuffer();
    let witness = new Uint8Array(witnessBuff);

    console.log("loading go env and wasm...");
    let wasi = new WASI({
        env: {},
        args: [
            `[${witness.join(",")}]`,
        ],
    });
    const wasmModule = await WebAssembly.compileStreaming(fetch(wasmUrl));
    const instance = await wasi.instantiate(wasmModule, {
        ...wasi.getImports(wasmModule)
    });
    console.log(instance)

    console.log("generating proof...");
    let status = wasi.start();
    console.log("status code", status);

    let res = wasi.getStdoutString();
    console.log("result", res);

    let end = Date.now();
    let elapsed = end - start;
    console.log("Finished!", `${ elapsed / 1000 }s`);
});