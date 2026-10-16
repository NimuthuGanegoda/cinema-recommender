package models

import "errors"

var (
	// ErrColomboExcluded is returned when a query targets the Colombo metropolitan district, which is excluded by policy.
	ErrColomboExcluded = errors.New("theaters in Colombo metropolitan district are excluded by policy (regional outstation theaters only)")

	// ErrCinemaNotFound is returned when a cinema identifier does not match any registered theater.
	ErrCinemaNotFound = errors.New("cinema not found in regional registry")

	// ErrEmptyCart is returned when evaluating an empty concession basket.
	ErrEmptyCart = errors.New("cart must contain at least one item")

	// ErrInvalidBudget is returned when the user specifies an invalid or non-positive budget.
	ErrInvalidBudget = errors.New("budget must be greater than zero")

	// ErrInvalidPartySize is returned when party size is less than 1.
	ErrInvalidPartySize = errors.New("party size must be at least 1")
)
