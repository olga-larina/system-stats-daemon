package model

import "errors"

var (
	ErrPeriodNotValid = errors.New("period not valid")
	ErrStatsNotValid  = errors.New("stats not valid")
)
