package repo

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

var ErrNoRows = "sql: no rows in result set"

// IsRecordNotFoundError returns current error has record not found error or not
func IsRecordNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	if err.Error() == ErrNoRows {
		return true
	}
	return gorm.IsRecordNotFoundError(err)
}

// IsRecordUniqueViolationError returns current error has record unique violation error or not
func IsRecordUniqueViolationError(err error) bool {
	if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
		return true
	}
	return false
}
