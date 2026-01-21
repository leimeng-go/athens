package mem

import (
	"fmt"

	"github.com/leimeng-go/athens/pkg/errors"
	"github.com/leimeng-go/athens/pkg/storage"
)

// NewStorage creates new in-memory storage for testing purposes.
// Note: This is a test-only implementation that uses MongoDB.
// For production, use mongo.NewStorage directly.
func NewStorage() (storage.Backend, error) {
	const op errors.Op = "mem.NewStorage"
	return nil, errors.E(op, fmt.Errorf("memory storage is not supported in this build - use MongoDB storage for testing"))
}
