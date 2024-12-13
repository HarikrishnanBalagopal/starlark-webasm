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

func convertToStarlarkValue(value js.Value) starlark.Value {
	switch value.Type() {
	case js.TypeBoolean:
		return starlark.Bool(value.Bool())
	case js.TypeNumber:
		floatVal := value.Float()
		if floatVal == float64(int(floatVal)) {
			return starlark.MakeInt(value.Int())
		}
		return starlark.Float(floatVal)
	case js.TypeString:
		return starlark.String(value.String())
	case js.TypeObject:
		if value.InstanceOf(js.Global().Get("Array")) {
			list := []starlark.Value{}
			length := value.Length()
			for i := 0; i < length; i++ {
				list = append(list, convertToStarlarkValue(value.Index(i)))
			}
			return starlark.NewList(list)
		} else {
			dict := starlark.NewDict(value.Length())
			keys := js.Global().Get("Object").Call("keys", value)
			length := keys.Length()
			for i := 0; i < length; i++ {
				key := keys.Index(i).String()
				dict.SetKey(starlark.String(key), convertToStarlarkValue(value.Get(key)))
			}
			return dict
		}
	default:
		return starlark.None
	}
}

func convertToJSValue(value starlark.Value) js.Value {
	switch v := value.(type) {
	case starlark.Bool:
		return js.ValueOf(bool(v))
	case starlark.Float:
		return js.ValueOf(float64(v))
	case starlark.String:
		return js.ValueOf(string(v))
	case starlark.Int:
		intVal, _ := v.Int64()
		return js.ValueOf(intVal)
	case *starlark.List:
		array := js.Global().Get("Array").New(v.Len())
		for i := 0; i < v.Len(); i++ {
			array.SetIndex(i, convertToJSValue(v.Index(i)))
		}
		return array
	case *starlark.Dict:
		obj := js.Global().Get("Object").New()
		for _, item := range v.Items() {
			key := item[0].(starlark.String)
			obj.Set(string(key), convertToJSValue(item[1]))
		}
		return obj
	default:
		return js.Null()
	}
}

func getStarlarkRunner() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 1 {
			err := fmt.Errorf("Error: expected at least one argument with the source code. Actual len(args) %d args %+v", len(args), args)
			return map[string]interface{}{"error": err.Error()}
		}
		starlark_code := args[0].String()
		funcName := "main"
		if len(args) > 1 {
			funcName = args[1].String()
		}
		funcArgs := []starlark.Value{}
		if len(args) > 2 {
			for _, arg := range args[2:] {
				funcArgs = append(funcArgs, convertToStarlarkValue(arg))
			}
		}

		output := strings.Builder{}
		thread := &starlark.Thread{Name: "js-go-starlark-thread", Print: func(_ *starlark.Thread, msg string) {
			output.WriteString(msg + "\n")
		}}
		globals, err := starlark.ExecFile(thread, "", starlark_code, nil)
		if err != nil {
			err := fmt.Errorf("Error: failed to evaluate the starlark code. Error: %q", err)
			return map[string]interface{}{"error": err.Error()}
		}
		mainFn, ok := globals[funcName]
		if !ok {
			err := fmt.Errorf("Error: the function %q is missing from the starlark code.", funcName)
			return map[string]interface{}{"error": err.Error()}
		}
		// Call the Starlark function from Go.
		result, err := starlark.Call(thread, mainFn, funcArgs, nil)
		if err != nil {
			err := fmt.Errorf("Error: failed to execute the starlark code. Error: %q", err)
			return map[string]interface{}{"error": err.Error()}
		}
		return map[string]interface{}{"message": output.String(), "returnValue": convertToJSValue(result)}
	})
}

func main() {
	js.Global().Set("run_starlark_code", getStarlarkRunner())
	fmt.Println("the run_starlark_code has been added to the javascript globals (window object)")
	<-make(chan bool) // keep thread running forever so Javascript can call the function we exported.
}
