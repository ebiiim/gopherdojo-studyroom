package conv

import "errors"

// Errors
var (
	ErrExtensionNotSupported = errors.New("ErrExtensionNotSupported")

	ErrCouldNotEncode = errors.New("ErrCouldNotEncode")
	ErrCouldNotDecode = errors.New("ErrCouldNotDecode")
	ErrCouldNotRead   = errors.New("ErrCouldNotRead")
	ErrCouldNotWrite  = errors.New("ErrCouldNotWrite")

	ErrFileAccess       = errors.New("ErrFileAccess")
	ErrFileNotSupported = errors.New("ErrFileNotSupported")
)
