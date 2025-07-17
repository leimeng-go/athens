package module

import (
	"context"

	"github.com/leimeng-go/athens/pkg/storage"
)

// UpstreamLister retrieves a list of available module versions from upstream
// i.e. VCS, and a Storage backend.
type UpstreamLister interface {
	List(ctx context.Context, mod string) (*storage.RevInfo, []string, error)
}
