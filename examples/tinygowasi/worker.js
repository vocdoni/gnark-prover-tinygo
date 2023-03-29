import './wasm_exec.js';

import { init, WASI } from '@wasmer/wasi';

import circuitCSUrl from "url:./artifacts/zkcensus.ccs";
import circuitSRSUrl from "url:./artifacts/zkcensus.srs";
import witnessUrl from "url:./artifacts/zkcensus.witness";
import wasmUrl from "url:./artifacts/circuit.wasm";

addEventListener('message', async(e) => {
    await init();

    let start = Date.now();
    console.log("reading artifacts...");

    let resCCS = await fetch(circuitCSUrl);
    let ccsBuff = await resCCS.arrayBuffer();
    let ccs = new Uint8Array(ccsBuff);
    console.log("(1/3)");

    let resSRS = await fetch(circuitSRSUrl);
    let srsBuff = await resSRS.arrayBuffer();
    let srs = new Uint8Array(srsBuff);
    console.log("(2/3)")

    let resWitness = await fetch(witnessUrl);
    let witnessBuff = await resWitness.arrayBuffer();
    let witness = new Uint8Array(witnessBuff);
    console.log("(3/3)");

    console.log("loading go env and wasm...");

    let wasi = new WASI({
        env: {},
        args: [
            "circuit.wasm",
            `[${ccs.join(",")}]`,
            `[${srs.join(",")}]`,
            `[${witness.join(",")}]`,
        ],
    });
    const wasmModule = await WebAssembly.compileStreaming(fetch(wasmUrl));
    const instance = await wasi.instantiate(wasmModule, {});
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