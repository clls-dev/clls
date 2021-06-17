package main

import (
	"fmt"
	"strings"
)

const (
	pkgName = "lspsrv"
	lspPkg  = "github.com/clls-dev/clls/pkg/lsp"
)

type methodDefinition struct {
	name           string
	id             string
	requestType    string
	responseType   string
	isNotification bool
}

var defs = []methodDefinition{
	// Notifications
	{"Initialized", "initialized", "lsp.InitializedParams", "", true},
	{"Exit", "exit", "", "", true},

	{"DidOpenTextDocument", "textDocument/didOpen", "lsp.DidOpenTextDocumentParams", "", true},
	{"DidCloseTextDocument", "textDocument/didClose", "lsp.DidCloseTextDocumentParams", "", true},
	{"DidChangeTextDocument", "textDocument/didChange", "lsp.DidChangeTextDocumentParams", "", true},
	{"DidSaveTextDocument", "textDocument/didSave", "lsp.DidSaveTextDocumentParams", "", true},

	// Calls
	{"Initialize", "initialize", "lsp.InitializeParams", "lsp.InitializeResult", false},
	{"Shutdown", "shutdown", "", "", false},

	{"SemanticTokens", "textDocument/semanticTokens/full", "lsp.SemanticTokensParams", "lsp.SemanticTokens", false},
	{"DocumentFormatting", "textDocument/formatting", "lsp.DocumentFormattingParams", "[]lsp.TextEdit", false},
	{"Rename", "textDocument/rename", "lsp.RenameParams", "lsp.WorkspaceEdit", false},
}

func (d *methodDefinition) fullRet() string {
	ret := "error"
	if d.responseType != "" && !d.isNotification {
		rt := d.responseType
		if !(strings.HasPrefix(rt, "[]") || strings.HasPrefix(rt, "map[")) {
			rt = "*" + rt
		}
		ret = fmt.Sprintf("(%s, error)", rt)
	}
	return ret
}

func (d *methodDefinition) fullReq() string {
	req := d.requestType
	if req == "" || strings.HasPrefix(req, "[]") {
		return req
	}
	return "*" + req
}
