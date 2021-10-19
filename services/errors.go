package services

import "errors"

var (
	ErrPaginationLimit = errors.New("pagination limit max is 100")
	ErrTryAgain        = errors.New("an error occured, please try again later")
)
