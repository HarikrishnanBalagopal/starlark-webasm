import { wasm } from '@rollup/plugin-wasm';

export default {
    input: 'src/index.js',
    output: {
        // dir: 'output',
        file: 'index.js',
        format: 'cjs'
    },
    plugins: [wasm()]
};