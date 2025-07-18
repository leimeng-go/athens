package gcp

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"github.com/leimeng-go/athens/pkg/config"
	"github.com/leimeng-go/athens/pkg/errors"
	"github.com/leimeng-go/athens/pkg/observ"
	pkgstorage "github.com/leimeng-go/athens/pkg/storage"
)

// Info implements Getter.
func (s *Storage) Info(ctx context.Context, module, version string) ([]byte, error) {
	const op errors.Op = "gcp.Info"
	ctx, span := observ.StartSpan(ctx, op.String())
	defer span.End()
	infoReader, err := s.bucket.Object(config.PackageVersionedName(module, version, "info")).NewReader(ctx)
	if err != nil {
		return nil, errors.E(op, err, getErrorKind(err), errors.M(module), errors.V(version))
	}
	infoBytes, err := io.ReadAll(infoReader)
	_ = infoReader.Close()
	if err != nil {
		return nil, errors.E(op, err, errors.M(module), errors.V(version))
	}
	return infoBytes, nil
}

// GoMod implements Getter.
func (s *Storage) GoMod(ctx context.Context, module, version string) ([]byte, error) {
	const op errors.Op = "gcp.GoMod"
	ctx, span := observ.StartSpan(ctx, op.String())
	defer span.End()
	modReader, err := s.bucket.Object(config.PackageVersionedName(module, version, "mod")).NewReader(ctx)
	if err != nil {
		return nil, errors.E(op, err, getErrorKind(err), errors.M(module), errors.V(version))
	}
	modBytes, err := io.ReadAll(modReader)
	_ = modReader.Close()
	if err != nil {
		return nil, errors.E(op, fmt.Errorf("could not get new reader for mod file: %w", err), errors.M(module), errors.V(version))
	}

	return modBytes, nil
}

// Zip implements Getter.
func (s *Storage) Zip(ctx context.Context, module, version string) (pkgstorage.SizeReadCloser, error) {
	const op errors.Op = "gcp.Zip"
	ctx, span := observ.StartSpan(ctx, op.String())
	defer span.End()
	zipReader, err := s.bucket.Object(config.PackageVersionedName(module, version, "zip")).NewReader(ctx)
	if err != nil {
		return nil, errors.E(op, err, getErrorKind(err), errors.M(module), errors.V(version))
	}
	return pkgstorage.NewSizer(zipReader, zipReader.Attrs.Size), nil
}

func getErrorKind(err error) int {
	if errors.IsErr(err, storage.ErrObjectNotExist) {
		return errors.KindNotFound
	}
	return errors.KindUnexpected
}
