importScripts("wasm_exec.js");

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

const go = new Go();

// Replace the console.log function in the Go environment
go.importObject.env["syscall/js.console_log"] = (sp) => {
  const s = go._inst.exports.ram.loadString(sp);
  postMessage({ type: "log", message: s });
};

onmessage = async (event) => {
  if (event.data.type === "generateProof") {
    const witness = event.data.witness;

    try {
      // Instantiate and run the wasm module with the modified Go environment
      const result = await WebAssembly.instantiateStreaming(fetch("/artifacts/g16_prover.wasm"), go.importObject);
      go.run(result.instance);

      // Call the generateProof function with the witness data
      const proof = generateProof(witness);

      // Send the result back to the main thread
      postMessage({ type: "proofGenerated", result: proof });
    } catch (error) {
      postMessage({ type: "error", message: error.message });
    }
  }
};
