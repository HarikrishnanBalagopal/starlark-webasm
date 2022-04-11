# Starlark interpreter for Javascript using WebAssembly

This library adds a function called `run_starlark_code` to the `window` object (globals).

## Usage

```console
$ npm install starlark-webasm
```

```js
import { initialize } from 'starlark-webasm';

const starlark_code = `
def main():
    print("hello world")
`;

async function main() {
    await initialize();
    const result = run_starlark_code(starlark_code);
    if(result.error) return console.error(result.error);
    console.log(result.message);
}

main();
```

Copy the WASM binary file `main.wasm` to the `dist/` folder.
```console
$ cp node_modules/starlark-webasm/main.wasm dist/
```
