package clls

var builtinFuncs = func() []*Function {
	names := []string{
		"if",
		"point_add", "c", "list",
		"l", "sha256", "f", "r", "pubkey_for_exp", "a", "x",
		"divmod", "substr", "concat", "logand", "qq", "unquote", "q",
		"quote", "i",
	}
	funcs := make([]*Function, len(names))
	for i, n := range names {
		funcs[i] = &Function{Name: &Token{Value: n}, Builtin: true}
	}
	return funcs
}()

var BuiltinFuncsByName = func() map[string]*Function {
	m := map[string]*Function{}
	for _, f := range builtinFuncs {
		m[f.Name.Value] = f
	}
	return m
}()
