import './wasm_exec.js';
import wasmFileName from 'url:./main.wasm';

async function initialize() {
    const go = new Go();
    const wasmModule = await WebAssembly.instantiateStreaming(fetch(wasmFileName), go.importObject);
    go.run(wasmModule.instance);
}

export { initialize };
