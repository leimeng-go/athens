package actions

import (
	"fmt"
	"time"

	"github.com/leimeng-go/athens/pkg/config"
	"github.com/leimeng-go/athens/pkg/errors"
	"github.com/leimeng-go/athens/pkg/storage"
	"github.com/leimeng-go/athens/pkg/storage/mongo"
)

// GetStorage returns storage backend based on env configuration.
// Only MongoDB storage is supported.
func GetStorage(storageType string, storageConfig *config.Storage, timeout time.Duration) (storage.Backend, error) {
	const op errors.Op = "actions.GetStorage"
	if storageType != "mongo" {
		return nil, errors.E(op, fmt.Sprintf("only 'mongo' storage type is supported, got: %s", storageType))
	}

	if storageConfig.Mongo == nil {
		return nil, errors.E(op, "Invalid Mongo Storage Configuration")
	}

	return mongo.NewStorage(storageConfig.Mongo, timeout)
}
