package main

import (
	"fmt"
	"reflect"
	"strings"

	lsp "go.lsp.dev/protocol"
)

func main() {
	var srv lsp.Server
	serverInterface := reflect.TypeOf(&srv).Elem()

	fmt.Print("package lspsrv\n\nimport (\n    lsp \"go.lsp.dev/protocol\"\n    \"encoding/json\"\n)\n\ntype UnimplementedLanguageServer struct {}\nvar _ lsp.Server = (*UnimplementedLanguageServer)(nil)\n")

	for i := 0; i < serverInterface.NumMethod(); i++ {
		method := serverInterface.Method(i)
		function := method.Type

		ins := []string(nil)
		for i := 0; i < function.NumIn(); i++ {
			inType := function.In(i)
			ins = append(ins, strings.Replace(inType.String(), "protocol.", "lsp.", -1))
		}

		outs := []string(nil)
		for i := 0; i < function.NumOut(); i++ {
			outType := function.Out(i)
			outs = append(outs, strings.Replace(outType.String(), "protocol.", "lsp.", -1))
		}
		outString := strings.Join(outs, ", ")
		if len(outs) > 1 {
			outString = "(" + outString + ")"
		}
		if outString != "" {
			outString = " " + outString
		}

		returnString := strings.Repeat("nil, ", len(outs)-1) + "ErrNotImplemented"
		if len(outs) == 1 {
			returnString = "nil"
		}

		fmt.Printf("\nfunc (*UnimplementedLanguageServer) %s(%s)%s {\n    return %s\n}\n", method.Name, strings.Join(ins, ", "), outString, returnString)
	}

	s := strings.TrimSpace(`
func Unmarshal(method string, payloadBytes []byte) (interface{}, error) {
	switch method {
`) + "\n"
	for i := 0; i < serverInterface.NumMethod(); i++ {
		method := serverInterface.Method(i)
		function := method.Type

		ins := []string(nil)
		for i := 0; i < function.NumIn(); i++ {
			inType := function.In(i)
			ins = append(ins, strings.Replace(inType.String(), "protocol.", "lsp.", -1))
		}

		methodID := "textDocument/" + uncap(method.Name)
		if specialID, ok := specialMethodIDs[method.Name]; ok {
			methodID = specialID
		}

		if method.Name == "Request" {
			// skip
		} else if len(ins) > 1 {
			paramsTypeString := strings.TrimPrefix(ins[1], "*")
			s += fmt.Sprintf(`case "%s":
			var payload %s
			return &payload, json.Unmarshal(payloadBytes, &payload)
		`, methodID, paramsTypeString) + "\n"
		} else {
			s += fmt.Sprintf(`case "%s":
				return nil, nil
			`, methodID) + "\n"
		}
	}
	s += strings.TrimSpace(`
	}
	return nil, ErrUnknownMethod
}
`) + "\n\n"
	fmt.Print(s)

	s = strings.TrimSpace(`
	func Request(ctx context.Context, s lsp.Server, method string, payload interface{}) (interface{}, error) {
		switch method {
	`) + "\n"
	for i := 0; i < serverInterface.NumMethod(); i++ {
		method := serverInterface.Method(i)
		function := method.Type

		ins := []string(nil)
		for i := 0; i < function.NumIn(); i++ {
			inType := function.In(i)
			ins = append(ins, strings.Replace(inType.String(), "protocol.", "lsp.", -1))
		}

		outs := []string(nil)
		for i := 0; i < function.NumOut(); i++ {
			outType := function.Out(i)
			outs = append(outs, strings.Replace(outType.String(), "protocol.", "lsp.", -1))
		}

		retPrefix := strings.Repeat("nil, ", 2-len(outs))

		methodID := "textDocument/" + uncap(method.Name)
		if specialID, ok := specialMethodIDs[method.Name]; ok {
			methodID = specialID
		}

		if method.Name == "Request" {
			// skip
		} else if len(ins) <= 1 {
			s += fmt.Sprintf(`case "%s":
					return %ss.%s(ctx)
				`, methodID, retPrefix, method.Name) + "\n"
		} else {
			s += fmt.Sprintf(`case "%s":
				castedPayload, ok := payload.(%s)
				if !ok {
					return nil, ErrBadPayloadType
				}
				return %ss.%s(ctx, castedPayload)
			`, methodID, ins[1], retPrefix, method.Name) + "\n"
		}
	}
	s += strings.TrimSpace(`
		}
		return nil, ErrUnknownMethod
	}
	`) + "\n\n"
	fmt.Print(s)

}

var specialMethodIDs = map[string]string{
	"Initialize":         "initialize",
	"Initialized":        "initialized",
	"Shutdown":           "shutdown",
	"Exit":               "exit",
	"SemanticTokensFull": "textDocument/semanticTokens/full",
}

func uncap(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[0:1]) + s[1:]
}
