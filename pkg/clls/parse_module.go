package clls

import (
	"fmt"
	"math/rand"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type module struct {
	Args         interface{}
	Constants    []*constant
	Functions    []*function
	funcsByName  map[string]*function
	constsByName map[string]*constant
	Main         *CodeBody
	Includes     map[string]*include
	ModToken     *Token
	IsMod        bool
	Comments     []*Token
}

func parseModules(l *zap.Logger, tree *ASTNode, readFile func(string) (string, error), tokens []*Token) ([]*module, error) {
	if tree == nil {
		return nil, errors.New("empty tree")
	}
	if len(tree.Children) == 0 {
		return nil, errors.New("no children in tree")
	}
	comments := []*Token(nil)
	for _, t := range tokens {
		if t.Kind == commentToken {
			comments = append(comments, t)
		}
	}
	mods := []*module(nil)
	for _, n := range tree.Children {
		n, ok := n.(*ASTNode)
		if !ok {
			continue
		}

		if len(n.Children) < 1 {
			continue
		}
		firstChild := n.Children[0]

		mod := &module{
			funcsByName:  map[string]*function{},
			constsByName: map[string]*constant{},
			Includes:     map[string]*include{},
			Comments:     comments,
		}

		if t, ok := firstChild.(*Token); ok {
			mod.ModToken = t
			if t.Value == "mod" {
				mod.IsMod = true
			}
		}

		children := n.Children
		if mod.IsMod && len(n.Children) > 2 {
			mod.Args = n.Children[1]
			children = n.Children[2:]
		}
		remaining := []*ASTNode(nil)
		for _, mn := range children {
			mn, ok := mn.(*ASTNode)
			if !ok {
				continue
			}

			if len(mn.Children) < 1 {
				remaining = append(remaining, mn)
				continue
			}
			t, ok := mn.Children[0].(*Token)
			if !ok {
				continue
			}
			switch t.Value {
			case "include":
				fincl := &include{
					Token: t,
				}
				var filePath string
				if len(mn.Children) > 1 {
					fincl.Value = mn.Children[1]
					if t, ok := mn.Children[1].(*Token); ok {
						filePath = t.Value
						var err error
						if fincl.Module, err = LoadCLVM(l, filePath, readFile); err != nil {
							fincl.Module = nil
							fincl.LoadError = err
						}
					}
				}
				if filePath == "" {
					filePath = "unknown-include-file-path-" + fmt.Sprint(rand.Uint64())
				}
				mod.Includes[filePath] = fincl
			case "defun", "defun-inline", "defmacro":
				f := &function{
					Raw:          mn,
					Inline:       t.Value == "defun-inline",
					Macro:        t.Value == "defmacro",
					KeywordToken: t,
				}
				mod.Functions = append(mod.Functions, f)

				var err error

				if len(mn.Children) > 1 {
					f.Name, ok = mn.Children[1].(*Token)
					if ok {
						mod.funcsByName[f.Name.Value] = f
					}
				}

				if len(mn.Children) > 2 {
					f.Params = mn.Children[2]
					if f.ParamsBody, err = parseBody(mod, mod.vars(), mn.Children[2]); err != nil {
						l.Error("failed to parse params body", zap.Error(err))
						f.ParamsBody = nil
					}
				}

				if len(mn.Children) > 3 {
					f.RawBody = mn.Children[3]
				}
			case "defconstant":
				c := &constant{
					Token: t,
				}
				mod.Constants = append(mod.Constants, c)

				if len(mn.Children) > 1 {
					c.Name = mn.Children[1]
					if t, ok := c.Name.(*Token); ok {
						mod.constsByName[t.Value] = c
					}
				}

				if len(mn.Children) > 2 {
					valBody, err := parseBody(mod, mod.vars(), mn.Children[2])
					if err == nil {
						c.Value = valBody
					}

				}
			default:
				remaining = append(remaining, mn)
			}
		}

		for _, f := range mod.Functions {
			vars := make(map[string]struct{})
			for k := range mod.vars() {
				vars[k] = struct{}{}
			}
			for k := range f.vars() {
				vars[k] = struct{}{}
			}
			var err error
			if f.Body, err = parseBody(mod, vars, f.RawBody); err != nil {
				l.Error("failed to parse function body", zap.Error(err))
				f.Body = nil
			}
		}

		if mod.IsMod {
			var err error
			if len(remaining) > 0 {
				if mod.Main, err = parseBody(mod, mod.vars(), remaining[len(remaining)-1]); err != nil {
					return nil, errors.Wrap(err, "parse body")
				}
			}
		}

		mods = append(mods, mod)
	}
	return mods, nil
}

type CodeBodyKind int

const (
	unknownBodyKind = CodeBodyKind(iota)
	IfBodyKind
	blockBodyKind
	CallBodyKind
	exceptionBodyKind
	OperatorBodyKind
	valueBodyKind
	ConstBodyKind
	VarBodyKind
	FuncVarBodyKind
)

type CodeBody struct {
	Raw      interface{}
	Kind     CodeBodyKind
	Children []*CodeBody
	Token    *Token

	Constant   *constant   `json:",omitempty"`
	Function   function    `json:",omitempty"`
	IfCond     *CodeBody   `json:",omitempty"`
	IfBranch   *CodeBody   `json:",omitempty"`
	ElseBranch *CodeBody   `json:",omitempty"`
	CallArgs   []*CodeBody `json:",omitempty"`
	opChildren []*CodeBody
}

func parseBody(mod *module, vars map[string]struct{}, tree interface{}) (*CodeBody, error) {
	if tree == nil {
		return nil, nil
	}
	switch tree := tree.(type) {
	case *Token:
		t := tree
		kind := valueBodyKind
		if t == nil {
			kind = blockBodyKind
		} else if _, ok := vars[t.Value]; ok {
			kind = VarBodyKind
		} else if _, ok := mod.constsByName[t.Value]; ok {
			kind = ConstBodyKind
		} else {
			for _, sm := range mod.Includes {
				if sm.Module == nil {
					continue
				}
				if _, ok := sm.Module.constsByName[t.Value]; ok {
					return &CodeBody{Kind: ConstBodyKind, Raw: tree, Token: t}, nil
				}
			}

			f, ok := mod.funcsByName[t.Text]
			if !ok {
				f, ok = builtinFuncsByName[t.Text]
			}
			if ok {
				return &CodeBody{Kind: FuncVarBodyKind, Raw: tree, Function: *f, Token: t}, nil
			}
		}
		return &CodeBody{Kind: kind, Raw: tree, Token: t}, nil
	case *ASTNode:
		if len(tree.Children) == 0 {
			return &CodeBody{Kind: blockBodyKind, Raw: tree}, nil
		}
		firstChildAsToken, ok := tree.Children[0].(*Token)
		if ok {
			switch firstChildAsToken.Kind {
			case basicToken:
				switch firstChildAsToken.Value {
				case "if":
					cb := &CodeBody{
						Raw:   tree,
						Kind:  IfBodyKind,
						Token: firstChildAsToken,
					}
					if len(tree.Children) > 1 {
						icb, err := parseBody(mod, vars, tree.Children[1])
						if err != nil {
							return nil, errors.Wrap(err, "parse if condition")
						}
						cb.IfCond = icb
					}
					if len(tree.Children) > 2 {
						ib, err := parseBody(mod, vars, tree.Children[2])
						if err != nil {
							return nil, errors.Wrap(err, "parse if branch")
						}
						cb.IfBranch = ib
					}
					if len(tree.Children) > 3 {
						eb, err := parseBody(mod, vars, tree.Children[3])
						if err != nil {
							return nil, errors.Wrap(err, "parse else branch")
						}
						cb.ElseBranch = eb
					}
					return cb, nil
				case "+", "-", "*", "/", ">", "=", ">s":
					ccb := []*CodeBody(nil)
					for _, c := range tree.Children[1:] {
						cb, err := parseBody(mod, vars, c)
						if err != nil {
							return nil, errors.Wrap(err, "parse operator child")
						}
						ccb = append(ccb, cb)
					}
					return &CodeBody{
						Raw:        tree,
						Kind:       OperatorBodyKind,
						opChildren: ccb,
						Children:   ccb,
					}, nil
				default:
					if t, ok := tree.Children[0].(*Token); ok {
						if c, ok := mod.constsByName[t.Value]; ok {
							return &CodeBody{
								Raw:      tree,
								Kind:     ConstBodyKind,
								Constant: c,
							}, nil
						}
						f, ok := mod.funcsByName[t.Value]
						if !ok {
							f, ok = builtinFuncsByName[t.Value]
						}
						if ok {
							args := []*CodeBody(nil)
							for _, e := range tree.Children[1:] {
								//fmt.Println("parsing code body", mod, vars, e)
								acb, err := parseBody(mod, vars, e)
								if err != nil {
									return nil, errors.Wrap(err, "parse call arg")
								}
								args = append(args, acb)
							}
							return &CodeBody{
								Raw:      tree,
								Kind:     CallBodyKind,
								Function: *f,
								CallArgs: args,
								Token:    t,
							}, nil
						}
					}
				}
			}
		}
		children := make([]*CodeBody, len(tree.Children))
		for i, c := range tree.Children {
			child, err := parseBody(mod, vars, c)
			if err != nil {
				return nil, errors.Wrap(err, "parse block body")
			}
			children[i] = child
		}
		return &CodeBody{
			Raw:      tree,
			Kind:     blockBodyKind,
			Children: children,
		}, nil
	default:
		return &CodeBody{
			Raw:  tree,
			Kind: unknownBodyKind,
		}, nil
	}
}
