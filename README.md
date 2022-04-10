# Starlark interpreter for Javascript using WebAssembly

This library adds a function called `run_starlark_code` to the `window` object (globals).

## Usage

```js
const starlark_code = `
def main():
    print("hello world")
`;
const result = run_starlark_code(starlark_code);
if(result.error) {
    console.error(result.error);
}else{
    console.log(result.message);
}
```
