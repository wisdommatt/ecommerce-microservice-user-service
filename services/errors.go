package services

import "errors"

var (
	errPaginationLimit = errors.New("pagination limit max is 100")
)
