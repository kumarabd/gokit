package inmem

import "github.com/kumarabd/gokit/errors"

var (
	ErrKeyNotExist = errors.New("", errors.Alert, "Key does not exist")
)
