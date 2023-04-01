importScripts('https://cdn.jsdelivr.net/gh/golang/go@go1.18.4/misc/wasm/wasm_exec.js')

addEventListener('message', async(e) => {
    let start = Date.now();
    console.log("loading go env and wasm...");
    const go = new Go();
    let wasmModule = await WebAssembly.instantiateStreaming(fetch("/wasm/circuit.wasm"), go.importObject);
    go.run(wasmModule.instance);

    console.log("reading artifacts...");

    let resCCS = await fetch("/artifacts/zkcensus.ccs");
    let ccsBuff = await resCCS.arrayBuffer();
    let ccs = new Uint8Array(ccsBuff);
    console.log("(1/4)")

    let resSRS = await fetch("/artifacts/zkcensus.srs");
    let srsBuff = await resSRS.arrayBuffer();
    let srs = new Uint8Array(srsBuff);
    console.log("(2/4)")

    let resPKey = await fetch("/artifacts/zkcensus.pkey");
    let pKeyBuff = await resPKey.arrayBuffer();
    let pKey = new Uint8Array(pKeyBuff);
    console.log("(3/4)")

    let resWitness = await fetch("/artifacts/zkcensus.witness");
    let witnessBuff = await resWitness.arrayBuffer();
    let witness = new Uint8Array(witnessBuff);
    console.log("(4/4)")

    console.log("generating proof...")
    console.log(generateProof(ccs, srs, pKey, witness));
    let end = Date.now();
    let elapsed = end - start;
    console.log("Finished!", `${elapsed / 1000}s`);
});