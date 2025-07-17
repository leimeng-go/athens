package mem

import (
	"testing"

	"github.com/leimeng-go/athens/pkg/index"
	"github.com/leimeng-go/athens/pkg/index/compliance"
)

func TestMem(t *testing.T) {
	indexer := &indexer{}
	compliance.RunTests(t, indexer, indexer.clear)
}

func (i *indexer) clear() error {
	i.mu.Lock()
	i.lines = []*index.Line{}
	i.mu.Unlock()
	return nil
}
