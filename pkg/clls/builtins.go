package clls

var builtinFuncs = func() []*function {
	names := []string{
		"if",
		"point_add", "c", "list",
		"l", "sha256", "f", "r", "pubkey_for_exp", "a", "x",
		"divmod", "substr", "concat", "logand", "qq", "unquote", "q",
		"quote",
	}
	funcs := make([]*function, len(names))
	for i, n := range names {
		funcs[i] = &function{Name: &Token{Value: n}, Builtin: true}
	}
	return funcs
}()

var builtinFuncsByName = func() map[string]*function {
	m := map[string]*function{}
	for _, f := range builtinFuncs {
		m[f.Name.Value] = f
	}
	return m
}()
