package main

import (
	"context"
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/clls-dev/clls/pkg/clls"
	"github.com/clls-dev/clls/pkg/examples"
	"go.uber.org/zap"
)

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("cllsPrettifyVHL", js.FuncOf(prettify))

	js.Global().Set("cllsPrettifyLisp", js.FuncOf(prettifyLisp))

	js.Global().Set("cllsInspect", js.FuncOf(inspect))

	js.Global().Set("cllsGetExamples", js.FuncOf(getExamples))

	println("Go WebAssembly Initialized")

	<-c
}

func prettify(a js.Value, b []js.Value) interface{} {
	cmd, r := clls.Command("")
	if err := cmd.ParseAndRun(context.Background(), []string{"--pp", b[0].String()}); err != nil {
		fmt.Printf("ERROR: %s\n\n", err)
		cmd.FlagSet.Usage()
	}
	return js.ValueOf(*r)
}

func prettifyLisp(a js.Value, b []js.Value) interface{} {
	cmd, r := clls.Command("")
	if err := cmd.ParseAndRun(context.Background(), []string{"--pp", b[0].String(), "--lisp"}); err != nil {
		fmt.Printf("ERROR: %s\n\n", err)
		cmd.FlagSet.Usage()
	}
	return js.ValueOf(*r)
}

func inspect(a js.Value, b []js.Value) interface{} {
	sources := map[string]string{"main.clvm": b[0].String()}
	exs, err := examples.F.ReadDir(".")
	if err == nil {
		for _, e := range exs {
			b, err := examples.F.ReadFile(e.Name())
			if err == nil {
				sources[e.Name()] = string(b)
			}
		}
	}

	l := zap.NewNop()

	mod, err := clls.LoadCLVMFromStrings(l, "main.clvm", "file://main.clvm", sources)
	if err != nil {
		return js.ValueOf(err.Error())
	}

	med, err := json.Marshal(mod)
	if err != nil {
		return js.ValueOf(err.Error())
	}

	var out interface{}
	if err := json.Unmarshal(med, &out); err != nil {
		return js.ValueOf(err.Error())
	}

	return js.ValueOf(out)
}

func getExamples(a js.Value, b []js.Value) interface{} {
	exs := make(map[string]interface{})
	entries, err := examples.F.ReadDir(".")
	if err != nil {
		return js.ValueOf(err.Error())
	}
	for _, e := range entries {
		b, err := examples.F.ReadFile(e.Name())
		if err == nil {
			exs[e.Name()] = string(b)
		}
	}
	return js.ValueOf(exs)
}
