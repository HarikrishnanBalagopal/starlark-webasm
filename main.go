// Copyright 2022 Harikrishnan Balagopal

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"strings"
	"syscall/js"

	"go.starlark.net/starlark"
)

func getStarlarkRunner() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			err := fmt.Errorf("Error: expected there to be exactly one argument with the source code. Actual len(args) %d args %+v", len(args), args)
			return map[string]interface{}{"error": err.Error()}
		}
		starlark_code := args[0].String()
		output := strings.Builder{}
		thread := &starlark.Thread{Name: "js-go-starlark-thread", Print: func(_ *starlark.Thread, msg string) {
			output.WriteString(msg + "\n")
		}}
		globals, err := starlark.ExecFile(thread, "", starlark_code, nil)
		if err != nil {
			err := fmt.Errorf("Error: failed to evaluate the starlark code. Error: %q", err)
			return map[string]interface{}{"error": err.Error()}
		}
		mainFn, ok := globals["main"]
		if !ok {
			err := fmt.Errorf("Error: the main function is missing from the starlark code.")
			return map[string]interface{}{"error": err.Error()}
		}
		// Call the Starlark main function from Go.
		if _, err := starlark.Call(thread, mainFn, nil, nil); err != nil {
			err := fmt.Errorf("Error: failed to execute the starlark code. Error: %q", err)
			return map[string]interface{}{"error": err.Error()}
		}
		return map[string]interface{}{"message": output.String()}
	})
}

func main() {
	js.Global().Set("run_starlark_code", getStarlarkRunner())
	fmt.Println("the run_starlark_code has been added to the javascript globals (window object)")
	<-make(chan bool) // keep thread running forever so Javascript can call the function we exported.
}
