package azureblob

import (
	"context"
	"fmt"
	"io"

	"github.com/leimeng-go/athens/pkg/config"
	"github.com/leimeng-go/athens/pkg/errors"
	"github.com/leimeng-go/athens/pkg/observ"
	"github.com/leimeng-go/athens/pkg/storage"
)

// Info implements the (./pkg/storage).Getter interface.
func (s *Storage) Info(ctx context.Context, module, version string) ([]byte, error) {
	const op errors.Op = "azureblob.Info"
	ctx, span := observ.StartSpan(ctx, op.String())
	defer span.End()

	exists, err := s.Exists(ctx, module, version)
	if err != nil {
		return nil, errors.E(op, err, errors.M(module), errors.V(version))
	}
	if !exists {
		return nil, errors.E(op, errors.M(module), errors.V(version), errors.KindNotFound)
	}

	infoReader, err := s.client.ReadBlob(ctx, config.PackageVersionedName(module, version, "info"))
	if err != nil {
		return nil, errors.E(op, err, errors.M(module), errors.V(version))
	}

	infoBytes, err := io.ReadAll(infoReader)
	if err != nil {
		return nil, errors.E(op, err, errors.M(module), errors.V(version))
	}

	err = infoReader.Close()
	if err != nil {
		return nil, errors.E(op, err, errors.M(module), errors.V(version))
	}

	return infoBytes, nil
}

// GoMod implements the (./pkg/storage).Getter interface.
func (s *Storage) GoMod(ctx context.Context, module, version string) ([]byte, error) {
	const op errors.Op = "azureblob.GoMod"
	ctx, span := observ.StartSpan(ctx, op.String())
	defer span.End()

	exists, err := s.Exists(ctx, module, version)
	if err != nil {
		return nil, errors.E(op, err, errors.M(module), errors.V(version))
	}
	if !exists {
		return nil, errors.E(op, errors.M(module), errors.V(version), errors.KindNotFound)
	}

	modReader, err := s.client.ReadBlob(ctx, config.PackageVersionedName(module, version, "mod"))
	if err != nil {
		return nil, errors.E(op, err, errors.M(module), errors.V(version))
	}

	modBytes, err := io.ReadAll(modReader)
	if err != nil {
		return nil, errors.E(op, fmt.Errorf("could not get new reader for mod file: %w", err), errors.M(module), errors.V(version))
	}

	err = modReader.Close()
	if err != nil {
		return nil, errors.E(op, err, errors.M(module), errors.V(version))
	}

	return modBytes, nil
}

// Zip implements the (./pkg/storage).Getter interface.
func (s *Storage) Zip(ctx context.Context, module, version string) (storage.SizeReadCloser, error) {
	const op errors.Op = "azureblob.Zip"
	ctx, span := observ.StartSpan(ctx, op.String())
	defer span.End()
	exists, err := s.Exists(ctx, module, version)
	if err != nil {
		return nil, errors.E(op, err, errors.M(module), errors.V(version))
	}
	if !exists {
		return nil, errors.E(op, errors.M(module), errors.V(version), errors.KindNotFound)
	}
	zipReader, err := s.client.ReadBlob(ctx, config.PackageVersionedName(module, version, "zip"))
	if err != nil {
		return nil, errors.E(op, err, errors.M(module), errors.V(version))
	}
	return zipReader, nil
}
