import './wasm_exec.js';
import mainWASM from './main.wasm';

async function initialize() {
    const go = new Go();
    const wasmModule = await mainWASM(go.importObject);
    // const wasmModule = await WebAssembly.instantiateStreaming(fetch(wasmFileName), go.importObject);
    go.run(wasmModule.instance);
}

export { initialize };
