package database

import (
	"errors"
)

var ErrNoRecord = errors.New("database: no matching record found")
