// Importing the latest version of @wasmer/wasi and @wasmer/wasmfs libraries.
//import { init, WASI } from "@wasmer/wasi";
//import { init, WASI } from 'https://deno.land/x/wasm/wasi.ts';
//import { WASI } from "https://deno.land/x/wasm@v1.2.2/pkg/wasmer_wasi_js.js";
//import { WasmFs } from "@wasmer/wasmfs";
//import { MemFS } from "https://deno.land/x/wasm@v1.2.2/wasi.ts";
//import { MemFS } from "https://deno.land/x/wasm@v1.2.2/pkg/wasmer_wasi_js.js";

//import {init, WASI} from "https://unpkg.com/@wasmer/wasi@1.2.2/dist/Library.umd.min.js";
//import WasmFs from '@wasmer/wasmfs';
import { init, WASI } from '@wasmer/wasi';

import './wasm_exec.js';

import witnessUrl from "url:./artifacts/zkcensus.witness";
import wasmUrl from "url:./artifacts/g16_prover.wasm";


//let module = undefined;

const workerConsole = {
    log: (message) => {
      postMessage({ type: "log", message: message });
    },
    error: (message) => {
      postMessage({ type: "error", message: message });
    },
  };

console.log = workerConsole.log;
console.error = workerConsole.error;


onmessage = async (event) => {
  if (event.data.type === "generateProof") {
  await init();
    
  try {
    console.log("reading artifacts...");

    let resWitness = await fetch(witnessUrl);
    let witnessBuff = await resWitness.arrayBuffer();
    let witness = new Uint8Array(witnessBuff);

    // Instantiate new WASI and WasmFs Instances
    //const wasmFs = new WasmFs();
    console.log("loading go env and wasm...");
    let wasi = new WASI({
      args: [  
        `[${witness.join(",")}]`,
      ],
      env: {},
      bindings: {
        ...WASI.defaultBindings,
      }
    });

      // Fetch our Wasm File
      console.log("fetching wasm...");
      let response = await fetch(wasmUrl)
      let wasmBytes = new Uint8Array(await response.arrayBuffer())

      // Instantiate the WebAssembly file
      console.log("instantiating wasm...");
      let wasmModule = await WebAssembly.compile(wasmBytes);
      let instance = await WebAssembly.instantiate(wasmModule, {
        ...wasi.getImports(wasmModule)
      });

      // Start the WebAssembly WASI instance!
      wasi.start(instance);
      let stdout = wasi.getStdoutString(); // Get the contents of stdout

      console.log(`Standard Output: ${stdout}`); // Write stdout data to the DOM

      // Send the result back to the main thread
      postMessage({ type: "proofGenerated", result: stdout });
    } catch (error) {
      postMessage({ type: "error", message: error.message });
    }
  }
};
