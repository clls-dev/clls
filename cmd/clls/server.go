package main

import (
	"io/ioutil"
	"strings"

	"github.com/clls-dev/clls/pkg/clls"
	"github.com/clls-dev/clls/pkg/lsp"
	"github.com/clls-dev/clls/pkg/lspsrv"
	"github.com/pkg/errors"
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

var _ lspsrv.LanguageServer = (*server)(nil)

func newServer(l *zap.Logger) *server {
	if l == nil {
		l = zap.NewNop()
	}
	return &server{
		l:          l,
		openedDocs: map[string]*documentData{},
		cache:      newDocumentCache(200),
	}
}

func (s *server) loadCLVM(uriStr string) (*clls.Module, error) {
	if d, ok := s.openedDocs[uriStr]; ok && d.parsedModule {
		return d.module, nil
	}

	mod, err := clls.LoadCLVM(s.l, uriStr, s.readFile)
	if err != nil {
		return nil, err
	}

	s.l.Debug("parse module", zap.String("uri", uriStr))

	if d, ok := s.openedDocs[uriStr]; ok {
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
		s.l.Debug("reading file from docs", zap.String("uri", uriStr), zap.Int("size", len(d.content)), zap.String("hash", shortString(d.contentHash, 7)))
		return d.content, nil
	}

	s.l.Debug("will read file", zap.String("uri", uriStr))

	if !strings.HasPrefix(uriStr, fileURIPrefix) {
		return "", errors.New("not a file uri")
	}
	pathStr := strings.TrimPrefix(uriStr, fileURIPrefix)

	b, err := ioutil.ReadFile(pathStr)
	if err != nil {
		return "", err
	}
	fileStr := string(b)

	s.l.Debug("did read file", zap.String("uri", uriStr))

	return fileStr, nil
}
