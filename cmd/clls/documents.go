package main

import (
	"encoding/base64"

	"github.com/clls-dev/clls/pkg/clls"
	"golang.org/x/crypto/sha3"
)

type documentData struct {
	content     string
	contentHash string

	parsedModule bool
	module       *clls.Module

	generatedTokens bool
	semanticTokens  []uint32

	generatedSymbols bool
	symbols          []*clls.Symbol
}

func hashString(s string) string {
	h := sha3.New256()
	if _, err := h.Write([]byte(s)); err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

func newDocumentData(text string) *documentData {
	return &documentData{
		content:     text,
		contentHash: hashString(text), // FIXME, if an include changes, the cache will be bad
	}
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
