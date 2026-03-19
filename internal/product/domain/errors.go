package domain

import sharedErrors "github.com/ramdhanrizki/bytecode-api/internal/shared/errors"

var (
	ErrProductNotFound = sharedErrors.NotFound("product not found")
)
