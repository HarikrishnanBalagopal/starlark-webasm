import './wasm_exec.js';

async function initialize(WASM_FILE_URL = 'main.wasm') {
    const go = new Go();
    const wasmModule = await WebAssembly.instantiateStreaming(fetch(WASM_FILE_URL), go.importObject);
    go.run(wasmModule.instance);
}

export { initialize };
