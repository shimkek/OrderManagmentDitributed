package common

import "fmt"

var (
	ErrNoItems       = fmt.Errorf("items cannot be empty")
	ErrOrderNotFound = fmt.Errorf("order not found")
	ErrNotInStock    = fmt.Errorf("items are not in stock")
)
