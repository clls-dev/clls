package clls

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/clls-dev/clls/pkg/examples"
	"github.com/clls-dev/clls/pkg/lsp"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	CommandName = "clls"
)

func readFileToString(p string) (string, error) {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Command(rootName string) (*ffcli.Command, *string) {
	flagSet := flag.NewFlagSet(fmt.Sprintf("%s %s", rootName, CommandName), flag.ExitOnError)

	ppFlag := flagSet.String("pp", "", `--pp "(mod () ("Hello world!"))"`)

	s := ""
	r := &s

	return &ffcli.Command{
		Name:       CommandName,
		ShortUsage: fmt.Sprintf("%s %s <arg>", rootName, CommandName),
		ShortHelp:  "parse clvm high level code",
		FlagSet:    flagSet,
		Exec: func(_ context.Context, args []string) error {
			l, err := zap.NewDevelopment()
			if err != nil {
				fmt.Println("failed to init zap logger", err)
				l = zap.NewNop()
			}

			var mod *Module
			if *ppFlag != "" {
				sources := map[string]string{"main.clvm": *ppFlag}
				exs, err := examples.F.ReadDir(".")
				if err == nil {
					for _, e := range exs {
						b, err := examples.F.ReadFile(e.Name())
						if err == nil {
							fmt.Println("source", e.Name())
							sources[e.Name()] = string(b)
						}
					}
				}

				nmod, err := LoadCLVMFromStrings(l, "main.clvm", "file://main.clvm", sources)
				if err != nil {
					return errors.Wrap(err, "parse clvm")
				}

				//fmt.Println(mod)
				mod = nmod
			} else {
				filePath := flagSet.Arg(0)
				if filePath == "" {
					return errors.New("missing file path argument")
				}
				nmod, err := LoadCLVM(l, "file://"+filePath, readFileToString)
				if err != nil {
					return errors.Wrap(err, "parse modules")
				}
				mod = nmod
			}
			//fmt.Println(mod)
			/*if *lispFlag {
				s = mod.lispString()
			} else {
				s = mod.String()
			}*/
			mbytes, err := json.MarshalIndent(mod, "", "    ")
			if err == nil {
				fmt.Println(string(mbytes))
			}
			return nil
		},
	}, r
}

const fileURIPrefix = "file://"

// TODO: replace p with documentURI
func LoadCLVM(l *zap.Logger, documentURI lsp.DocumentURI, readFile func(string) (string, error)) (*Module, error) {
	f, err := readFile(documentURI)
	if err != nil {
		return nil, errors.Wrap(err, "read file")
	}

	tokens := []*Token(nil)
	tch, errptr := tokenize(f, documentURI)
	duptch := make(chan *Token)
	go func() {
		defer close(duptch)
		for token := range tch {
			duptch <- token
			tokens = append(tokens, token)
		}
	}()

	ast, err := parseAST(duptch)
	if err != nil {
		return nil, errors.Wrap(err, "parse syntax tree")
	}

	if *errptr != nil {
		return nil, errors.Wrap(*errptr, "tokenize")
	}

	mods, err := parseModules(l, ast, documentURI, readFile, tokens)
	if err != nil {
		return nil, errors.Wrap(err, "parse modules")
	}
	if len(mods) == 0 {
		return nil, errors.New("no modules in file")
	}

	return mods[0], nil

}

func LoadCLVMFromStrings(l *zap.Logger, p string, documentURI lsp.DocumentURI, files map[lsp.DocumentURI]string) (*Module, error) {
	return LoadCLVM(l, documentURI, func(p lsp.DocumentURI) (string, error) {
		f, ok := files[p]
		if !ok {
			return "", fmt.Errorf("unknown file '%s'", p)
		}
		return f, nil
	})
}

type constant struct {
	Token *Token
	Name  interface{}
	Value *CodeBody
}

type Function struct {
	Raw          *ASTNode
	RawBody      interface{}
	KeywordToken *Token
	Name         *Token
	Params       interface{}
	ParamsBody   *CodeBody
	Body         *CodeBody
	Inline       bool `json:",omitempty"`
	Builtin      bool `json:",omitempty"`
	Macro        bool `json:",omitempty"`
}

func paramsNames(n interface{}) map[string]struct{} {
	toks := paramsTokens(n)
	v := make(map[string]struct{}, len(toks))
	for _, t := range toks {
		v[t.Value] = struct{}{}
	}
	return v
}

func paramsTokens(n interface{}) map[string]*Token {
	v := map[string]*Token{}
	switch n := n.(type) {
	case *Token:
		if n.Kind == basicToken && n.Text != "." {
			v[n.Text] = n
		}
	case *ASTNode:
		for _, c := range n.Children {
			for k, sv := range paramsTokens(c) {
				v[k] = sv
			}
		}
	}
	return v
}

func (f *Function) vars() map[string]struct{} {
	return paramsNames(f.Params)
}

func (f *Function) varTokens() map[string]*Token {
	return paramsTokens(f.Params)
}

type include struct {
	Token     *Token
	Value     interface{}
	Module    *Module
	LoadError error
}

func (m *Module) vars() map[string]struct{} {
	return paramsNames(m.Args)
}

func (m *Module) varTokens() map[string]*Token {
	return paramsTokens(m.Args)
}

var ConditionCodes = func() string {
	b, err := examples.F.ReadFile("condition_codes.clvm")
	if err == nil {
		return string(b)
	}
	return ""
}()
