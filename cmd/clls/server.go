package main

import (
	"io/ioutil"
	"path/filepath"

	"github.com/clls-dev/clls/pkg/clls"
	"github.com/clls-dev/clls/pkg/lspsrv"
	"go.uber.org/zap"
)

type server struct {
	lspsrv.UnimplementedLanguageServer

	down bool
	exit bool
	docs map[string]string

	l *zap.Logger
}

func newServer(l *zap.Logger) *server {
	if l == nil {
		l = zap.NewNop()
	}
	return &server{
		l:    l,
		docs: map[string]string{},
	}
}

func (s *server) loadCLVM(filename, dirname, uriStr string) (*clls.Module, error) {
	return clls.LoadCLVM(s.l, filename, uriStr, func(p string) (string, error) {
		if d, ok := s.docs[fileURIPrefix+filepath.Join(dirname, p)]; ok {
			return d, nil
		}
		b, err := ioutil.ReadFile(filepath.Join(dirname, p))
		if err != nil {
			return "", err
		}
		return string(b), err
	})
}
