package clls

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/clls-dev/clls/pkg/examples"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	CommandName = "clls"
)

func Command(rootName string) (*ffcli.Command, *string) {
	flagSet := flag.NewFlagSet(fmt.Sprintf("%s %s", rootName, CommandName), flag.ExitOnError)

	ppFlag := flagSet.String("pp", "", `--pp "(mod () ("Hello world!"))"`)
	//lispFlag := flagSet.Bool("lisp", false, "--lisp: outputs clvm")

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

			var mod *module
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

				nmod, err := LoadCLVMFromStrings(l, "main.clvm", sources)
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
				nmod, err := LoadCLVM(l, filePath, func(p string) (string, error) {
					b, err := ioutil.ReadFile(p)
					if err != nil {
						return "", err
					}
					return string(b), nil
				})
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

func LoadCLVM(l *zap.Logger, p string, readFile func(string) (string, error)) (*module, error) {
	f, err := readFile(p)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("read file '%s'", p))
	}

	tch, ech := tokenize(f)

	tokens := []*Token(nil)
	duptch := make(chan *Token)

	dupech := make(chan error)
	go func() {
		defer close(dupech)
		dupech <- func() error {
			defer close(duptch)
			for {
				select {
				case token, open := <-tch:
					if !open {
						return nil
					}
					duptch <- token
					tokens = append(tokens, token)
				case err := <-ech:
					return err
				}
			}
		}()
	}()

	ast, err := parseAST(duptch)
	if err != nil {
		return nil, errors.Wrap(err, "parse syntax tree")
	}

	if err := <-dupech; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("tokenize file '%s'", p))
	}

	mods, err := parseModules(l, ast, readFile, tokens)
	if err != nil {
		return nil, errors.Wrap(err, "parse modules")
	}
	if len(mods) == 0 {
		return nil, errors.New("no modules in file")
	}

	return mods[0], nil

}

func LoadCLVMFromStrings(l *zap.Logger, p string, files map[string]string) (*module, error) {
	return LoadCLVM(l, p, func(p string) (string, error) {
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

type function struct {
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
	v := map[string]struct{}{}
	switch n := n.(type) {
	case *Token:
		v[n.Value] = struct{}{}
	case *ASTNode:
		for _, c := range n.Children {
			cn := paramsNames(c)
			for k := range cn {
				if k == "." {
					continue
				}
				v[k] = struct{}{}
			}
		}
	}
	return v
}

func (f *function) vars() map[string]struct{} {
	return paramsNames(f.Params)
}

type include struct {
	Token     *Token
	Value     interface{}
	Module    *module
	LoadError error
}

func (m *module) vars() map[string]struct{} {
	return paramsNames(m.Args)
}

var ConditionCodes = func() string {
	b, err := examples.F.ReadFile("condition_codes.clvm")
	if err == nil {
		return string(b)
	}
	return ""
}()
