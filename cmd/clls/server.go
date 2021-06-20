package main

import (
	"context"
	"io/ioutil"

	"github.com/clls-dev/clls/pkg/clls"
	"github.com/clls-dev/clls/pkg/lspsrv"
	lsp "go.lsp.dev/protocol"
	"go.uber.org/zap"
)

type server struct {
	lspsrv.UnimplementedLanguageServer

	down bool
	exit bool

	openedDocs map[lsp.DocumentURI]*documentData
	cache      *documentCache

	l *zap.Logger
}

var _ lsp.Server = (*server)(nil)

func newServer(l *zap.Logger) *server {
	if l == nil {
		l = zap.NewNop()
	}
	return &server{
		l:          l,
		openedDocs: map[lsp.DocumentURI]*documentData{},
		cache:      newDocumentCache(200),
	}
}

func (s *server) loadCLVM(u lsp.DocumentURI) (*clls.Module, error) {
	if d, ok := s.openedDocs[u]; ok && d.parsedModule {
		return d.module, nil
	}

	mod, err := clls.LoadCLVM(s.l, u, s.readFile)
	if err != nil {
		return nil, err
	}

	s.l.Debug("parse module", zap.Any("uri", u))

	if d, ok := s.openedDocs[u]; ok {
		d.module = mod
		d.parsedModule = true
	}
	return mod, nil
}

func shortString(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func (s *server) readFile(uriStr lsp.DocumentURI) (string, error) {
	if d, ok := s.openedDocs[uriStr]; ok {
		s.l.Debug("reading file from docs", zap.Any("uri", uriStr), zap.Int("size", len(d.content)), zap.String("hash", shortString(d.contentHash, 7)))
		return d.content, nil
	}

	s.l.Debug("will read file", zap.Any("uri", uriStr))

	b, err := ioutil.ReadFile(uriStr.Filename())
	if err != nil {
		return "", err
	}
	fileStr := string(b)

	s.l.Debug("did read file", zap.Any("uri", uriStr))

	return fileStr, nil
}

func (s *server) Request(ctx context.Context, method string, params interface{}) (interface{}, error) {
	return lspsrv.Request(ctx, s, method, params)
}
