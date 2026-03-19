package domain

import sharedErrors "github.com/ramdhanrizki/bytecode-api/internal/shared/errors"

var (
	ErrCategoryNotFound = sharedErrors.NotFound("category not found")
)
