package models

import "errors"

var (
	// ErrColomboExcluded is returned when a query targets the Colombo metropolitan district, which is excluded by policy.
	ErrColomboExcluded = errors.New("theaters in Colombo metropolitan district are excluded by policy (regional outstation theaters only)")

	// ErrCinemaNotFound is returned when a cinema identifier does not match any registered theater.
	ErrCinemaNotFound = errors.New("cinema not found in regional registry")

	// ErrEmptyCart is returned when evaluating or checking out an empty concession basket.
	ErrEmptyCart = errors.New("cart must contain at least one item")

	// ErrInvalidBudget is returned when the user specifies an invalid or non-positive budget.
	ErrInvalidBudget = errors.New("budget must be greater than zero")

	// ErrInvalidPartySize is returned when party size is less than 1.
	ErrInvalidPartySize = errors.New("party size must be at least 1")

	// ErrPaymentMethodNotSupported is returned when an unrecognized or unsupported payment method is requested.
	ErrPaymentMethodNotSupported = errors.New("requested payment method is not supported")
)
