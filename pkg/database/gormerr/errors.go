package gormerr

import (
	"errors"

	"gorm.io/gorm"
)

// IsRecordNotFound returns true when the error indicates that the requested
// record does not exist in the database.
func IsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// IsDuplicateKey returns true when the error indicates a unique constraint
// violation (postgres error code 23505). It recognizes the GORM translated
// error, which is produced because the connection is opened with
// TranslateError enabled.
func IsDuplicateKey(err error) bool {
	return err != nil && errors.Is(err, gorm.ErrDuplicatedKey)
}
