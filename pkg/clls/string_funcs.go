package clls

import (
	"fmt"
	"strings"
)

/*

func (n *ASTNode) lispString() string {
	if n == nil {
		return "nil"
	}

	if len(n.Children) == 0 {
		if n.Token == nil {
			return "()"
		}
		if n.Token.Kind == quoteToken {
			return n.Token.Text
		}
		return n.Token.Value
	}
	strs := make([]string, len(n.Children))
	for i, c := range n.Children {
		strs[i] = c.String()
	}

	return "(" + strings.Join(strs, " ") + ")"
}

func (n *ASTNode) String() string {
	if n == nil {
		return "nil"
	}

	if len(n.Children) == 0 {
		if n.Token == nil {
			return "()"
		}
		if n.Token.Kind == quoteToken {
			return n.Token.Text
		}
		return n.Token.Value
	}
	strs := make([]string, len(n.Children))
	for i, c := range n.Children {
		strs[i] = c.String()
	}

	return "(" + strings.Join(strs, " ") + ")"
}

func argsListString(list *ASTNode, sep string, transform func(string) string) string {
	if transform == nil {
		transform = func(s string) string { return s }
	}
	if len(list.Children) == 0 {
		if list.Token == nil {
			return ""
		}
		return "..." + transform(list.Token.Value)
	}
	strs := []string(nil)
	for i := 0; i < len(list.Children); i++ {
		n := list.Children[i]
		if n.String() != "." || i >= len(list.Children)-1 {
			if len(n.Children) != 0 {
				strs = append(strs, "("+argsListString(n, ", ", transform)+")")
			} else {
				strs = append(strs, transform(n.String()))
			}
			continue
		}
		strs = append(strs, "..."+transform(list.Children[i+1].String()))
		i++

	}
	return strings.Join(strs, sep)
}

func argsListLispString(list *ASTNode, sep string, transform func(string) string) string {
	if transform == nil {
		transform = func(s string) string { return s }
	}
	if len(list.Children) == 0 {
		if list.Token == nil {
			return "()"
		}
		return transform(list.Token.Value)
	}
	strs := []string(nil)
	for i := 0; i < len(list.Children); i++ {
		n := list.Children[i]
		if len(n.Children) != 0 {
			strs = append(strs, argsListLispString(n, sep, transform))
		} else {
			strs = append(strs, transform(n.lispString()))
		}
		i++

	}
	return "(" + strings.Join(strs, sep) + ")"
}

func (m *module) lispString() string {
	if m == nil {
		return "(*module)(nil)"
	}

	outer := "("
	if m.IsMod {
		outer += colorString(100, 100, 255, "mod") + " " + argsListLispString(m.Args, " ", colorArgs) + " "
	}
	outer += "\n"

	s := ""

	for _, c := range m.Includes {
		s += fmt.Sprintf("(%s %s)\n", colorString(100, 100, 255, "include"), c)
	}
	if len(m.Includes) != 0 {
		s += "\n"
	}
	for _, c := range m.Constants {
		s += fmt.Sprintf("(%s %s %s)\n", colorString(100, 100, 255, "defconstant"), constColor(c.Name.Token.Value), c.Value)
	}
	if len(m.Constants) != 0 && (len(m.Functions) != 0 || m.Main != nil) {
		s += "\n"
	}
	for _, f := range m.Functions {
		tname := "defun"
		if f.Inline {
			tname = "defun-inline"
		}
		if f.Macro {
			tname = "defmacro"
		}
		s += fmt.Sprintf("(%s %s %s \n%s\n)\n\n", colorString(100, 100, 255, tname), colorString(255, 255, 0, f.Name.Token.Value), argsListLispString(f.Params, " ", colorArgs), prependIndents(f.Body.lispString(), indent))
	}

	outersuffix := ""
	if m.Main != nil {
		s += "(" + strings.TrimSpace(m.Main.lispString()) + ")"
		outersuffix += "\n"
	}
	outersuffix += ")"

	return outer + prependIndents(s, indent) + outersuffix
}

func (m *module) String() string {
	if m == nil {
		return "(*module)(nil)"
	}

	outer := ""
	if m.IsMod {
		outer += colorString(100, 100, 255, "mod") + " "
	}
	outer += "{\n"

	s := ""

	for _, c := range m.Includes {
		s += fmt.Sprintf(`%s "%s"`+"\n", colorString(100, 100, 255, "include"), c)
	}
	if len(m.Includes) != 0 {
		s += "\n"
	}
	for _, c := range m.Constants {
		s += fmt.Sprintf("%s %s = %s\n", colorString(100, 100, 255, "const"), constColor(c.Name.Token.Value), c.Value)
	}
	if len(m.Constants) != 0 && (len(m.Functions) != 0 || m.Main != nil) {
		s += "\n"
	}
	for _, f := range m.Functions {
		prefix := ""
		if f.Inline {
			prefix += "inline-"
		}
		tname := "func"
		if f.Macro {
			tname = "macro"
		}
		s += fmt.Sprintf("%s %s(%s) {\n%s\n}\n\n", colorString(100, 100, 255, prefix+tname), colorString(255, 255, 0, f.Name.Token.Value), argsListString(f.Params, ", ", colorArgs), prependIndents(f.Body.String(), indent))
	}

	outersuffix := ""
	if m.Main != nil {
		s += colorString(100, 100, 255, "main-func") + " (" + argsListString(m.Args, ", ", colorArgs) + ") {\n" + prependIndents(m.Main.String(), indent) + "\n}"
		outersuffix += "\n"
	}
	outersuffix += "}"

	return outer + prependIndents(s, indent) + outersuffix
}

func (cb *CodeBody) String() string {
	if cb == nil {
		return "(*CodeBody)(nil)"
	}
	switch cb.Kind {
	case IfBodyKind:
		n := cb
		s := ""
		for n.Kind == IfBodyKind {
			s += fmt.Sprintf("%s %s {\n%s\n} %s ", colorString(255, 0, 255, "if"), n.IfCond.String(), prependIndents(n.IfBranch.String(), indent), colorString(255, 0, 255, "else"))
			n = n.ElseBranch
		}
		s += fmt.Sprintf("{\n%s\n}", prependIndents(n.String(), indent))
		return s
	case callBodyKind:
		if len(cb.Raw.Children) == 0 {
			return colorString(255, 255, 0, cb.Raw.Token.Value)
		}

		strs := []string(nil)
		for _, a := range cb.CallArgs {
			strs = append(strs, a.String())
		}
		prefix := colorString(255, 255, 0, cb.Function.Name.Token.Value) + "("
		suffix := ")"
		if cb.Function.Name.Token.Value == "list" && cb.Function.Builtin {
			prefix = "["
			suffix = "]"
		} else if cb.Function.Name.Token.Value == "c" && cb.Function.Builtin {
			s := cb.CallArgs[1].String()
			if cb.CallArgs[1].Kind == IfBodyKind {
				s = "(" + s + ")"
			}
			return fmt.Sprintf("[%s, ...%s]", cb.CallArgs[0], s)
		} else if cb.Function.Name.Token.Value == "x" && cb.Function.Builtin {
			prefix = colorString(255, 0, 0, "throw") + " ("
		} else if cb.Function.Name.Token.Value == "f" && cb.Function.Builtin {
			n := cb.CallArgs[0]
			d := 0
			for n.Kind == callBodyKind && n.Function.Name.Token.Value == "r" {
				n = n.CallArgs[0]
				d++
			}
			return fmt.Sprintf("%s[%d]", n.String(), d)
		} else if cb.Function.Name.Token.Value == "r" && cb.Function.Builtin {
			n := cb.CallArgs[0]
			d := 1
			for n.Kind == callBodyKind && n.Function.Name.Token.Value == "r" {
				n = n.CallArgs[0]
				d++
			}
			return fmt.Sprintf("%s[%d:]", n.String(), d)
		} else if cb.Function.Name.Token.Value == "qq" && cb.Function.Builtin {
			prefix = colorString(255, 150, 75, "```")
			suffix = colorString(255, 150, 75, "```")
			inner := strings.Join(strs, ", ")
			if strings.ContainsRune(inner, '\n') {
				prefix += "\n"
				suffix = "\n" + suffix
			}
			return prefix + inner + suffix
		} else if cb.Function.Name.Token.Value == "unquote" && cb.Function.Builtin {
			prefix = colorString(255, 150, 75, "${")
			suffix = colorString(255, 150, 75, "}")
		} else if cb.Function.Name.Token.Value == "q" && cb.Function.Builtin && len(cb.CallArgs) > 1 && cb.CallArgs[0].Raw.Token.Value == "." {
			prefix = colorString(255, 255, 0, cb.Function.Name.Token.Value) + "(..."
			return prefix + cb.CallArgs[1].String() + ")"
		}
		return prefix + strings.Join(strs, ", ") + suffix
	case constBodyKind:
		return constColor(cb.Raw.String())
	case varBodyKind:
		return colorString(160, 255, 160, cb.Raw.String())
	case OperatorBodyKind:
		strs := []string(nil)
		for _, a := range cb.opChildren {
			strs = append(strs, a.String())
		}
		v := cb.Raw.Children[0].Token.Value
		if v == "=" {
			v = "=="
		}
		return strings.Join(strs, " "+v+" ")
	}
	return cb.Raw.String()
}

func (cb *CodeBody) lispString() string {
	if cb == nil {
		return "(*CodeBody)(nil)"
	}
	switch cb.Kind {
	case IfBodyKind:
		n := cb
		return fmt.Sprintf("(%s %s\n%s\n%s\n)", colorString(255, 0, 255, "if"), n.IfCond.lispString(), prependIndents(n.IfBranch.lispString(), indent), prependIndents(n.ElseBranch.lispString(), indent))
	case callBodyKind:
		if len(cb.Raw.Children) == 0 {
			return colorString(255, 255, 0, cb.Raw.Token.Value)
		}

		strs := []string(nil)
		for _, a := range cb.CallArgs {
			strs = append(strs, a.lispString())
		}
		prefix := "(" + colorString(255, 255, 0, cb.Function.Name.Token.Value) + " "
		suffix := ")"
		if cb.Function.Name.Token.Value == "x" && cb.Function.Builtin {
			prefix = "(" + colorString(255, 0, 0, "x")
			if len(cb.Raw.Children) > 1 {
				prefix += " "
			}
		} else if cb.Function.Name.Token.Value == "qq" && cb.Function.Builtin {
			prefix = "(" + colorString(255, 150, 75, "qq") + " "
			suffix = ")"
			inner := strings.TrimSpace(strings.Join(strs, " "))
			return prefix + inner + suffix
		} else if cb.Function.Name.Token.Value == "unquote" && cb.Function.Builtin {
			prefix = "(" + colorString(255, 150, 75, "unquote") + " "
			suffix = ")"
		}
		return prefix + strings.Join(strs, " ") + suffix
	case constBodyKind:
		return constColor(cb.Raw.lispString())
	case varBodyKind:
		return colorString(160, 255, 160, cb.Raw.lispString())
	case OperatorBodyKind:
		strs := []string(nil)
		for _, a := range cb.opChildren {
			strs = append(strs, a.lispString())
		}
		v := cb.Raw.Children[0].Token.Value
		return "(" + v + " " + strings.Join(strs, " ") + ")"
	}
	return cb.Raw.lispString()
}

func colorArgs(s string) string {
	return colorString(160, 255, 160, s)
}

func constColor(s string) string {
	return colorString(150, 150, 255, s)
}

func prependIndents(s string, count int) string {
	strs := strings.Split(s, "\n")
	for i := range strs {
		strs[i] = strings.Repeat(" ", count) + strs[i]
	}
	return strings.Join(strs, "\n")
}

func colorString(r int, g int, b int, s string) string {
	return fmt.Sprintf(`<span style="color: rgb(%d, %d, %d)">%s</span>`, r, g, b, s)
	//return fmt.Sprintf("\033[38;2;%d;%d;%dm%s\033[00m", r, g, b, s)
}

const indent = 4*/

func (bk CodeBodyKind) String() string {
	s, ok := map[CodeBodyKind]string{
		IfBodyKind:        "if",
		blockBodyKind:     "block",
		CallBodyKind:      "call",
		exceptionBodyKind: "exception",
		OperatorBodyKind:  "op",
		valueBodyKind:     "value",
		ConstBodyKind:     "const",
		VarBodyKind:       "var",
		FuncVarBodyKind:   "func-var",
	}[bk]
	if !ok {
		return fmt.Sprintf("unknown(%d)", bk)
	}
	return s
}

func (t Token) String() string {
	if t.Kind == parensOpenToken || t.Kind == parensCloseToken || t.Kind == spaceToken {
		return t.Kind.String()
	}
	return fmt.Sprintf("%s(%s)", t.Kind, strings.TrimSpace(t.Value))
}

func (tk tokenKind) String() string {
	s, ok := tokenKindNames[tk]
	if !ok {
		return fmt.Sprintf("unknown(%d)", tk)
	}
	return s
}
