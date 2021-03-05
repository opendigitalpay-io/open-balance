package storage

import (
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type NotFoundError struct {
	Field string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s not found.", e.Field)
}

type DuplicatedEntryError struct {
	Field string
}

func (e DuplicatedEntryError) Error() string {
	return fmt.Sprintf("%s already exists.", e.Field)
}

func wrapDBError(err error, field string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NotFoundError{Field: field}
	}

	var mse *mysql.MySQLError
	if errors.As(err, &mse) {
		switch mse.Number {
		case 1062:
			return DuplicatedEntryError{Field: field}
		}
	}

	return err
}
