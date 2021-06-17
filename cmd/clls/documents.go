package main

import (
	"github.com/clls-dev/clls/pkg/clls"
	"github.com/clls-dev/clls/pkg/lsp"
)

type documentData struct {
	content     string
	contentHash string

	parsedModule bool
	module       *clls.Module

	generatedTokens bool
	semanticTokens  []lsp.UInteger

	generatedSymbols bool
	symbols          []*clls.Symbol
}

type documentCacheEntry struct {
	data  *documentData
	index int
}

type documentCache struct {
	size   int
	byHash map[string]*documentCacheEntry
	sorted []*documentCacheEntry
}

func newDocumentCache(size int) *documentCache {
	return &documentCache{
		size:   size,
		byHash: map[string]*documentCacheEntry{},
	}
}

func (dc *documentCache) put(dd *documentData) bool {
	if _, ok := dc.byHash[dd.contentHash]; ok {
		return false
	}
	if len(dc.sorted) >= dc.size {
		dc.sorted = dc.sorted[:dc.size]
	}
	entry := &documentCacheEntry{data: dd}
	dc.sorted = append([]*documentCacheEntry{entry}, dc.sorted...)
	for i := 1; i < len(dc.sorted); i++ {
		dc.sorted[i].index = i
	}
	dc.byHash[dd.contentHash] = entry
	return true
}

func (dc *documentCache) pull(contentHash string) (*documentData, bool) {
	ce, ok := dc.byHash[contentHash]
	if !ok {
		return nil, false
	}
	delete(dc.byHash, contentHash)
	dc.sorted = append(dc.sorted[:ce.index], dc.sorted[ce.index+1:]...)
	for i := ce.index; i < len(dc.sorted); i++ {
		dc.sorted[i].index = i
	}
	return ce.data, true
}
