# Starlark interpreter for Javascript using WebAssembly

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

## Overview

This library adds a function called `run_starlark_code` to the `window` object (globals).

Starlark is a Python-like language.  
It runs in a sandbox, meaning it does not have access to the file system, the network, the current time, sources of randomness, etc.
This makes it very safe to run unsafe user provided code because the Starlark code can only access stuff that we explicitly give it access to.

[Starlark](https://github.com/bazelbuild/starlark)  An overview of the language.  
[Starlark Language Specification](https://github.com/bazelbuild/starlark/blob/master/spec.md)  Detailed reference spec.  
[starlark-go](https://github.com/google/starlark-go) A Starlark interpreter written in Golang.

This interpreter and some wrapper code (see `main.go`) are compiled into a WASM binary.  
The javascript code in `src/index.js` exports a function called `initialize` that fetches the WASM module at runtime and compiles it.  
The WASM code adds a javascript function called `run_starlark_code` which accepts Starlark source code and returns an object.  
The object returned by `run_starlark_code` will have a field called `error` containing an error message OR  
it will have a field called `message` which contains the result of running the Starlark code.  
The output of all the `print` function calls in the Starlark code is returned as the result.
