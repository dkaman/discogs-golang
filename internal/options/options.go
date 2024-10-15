package options

import (
	"fmt"
)

type Option[T any] func(*T) error

func Apply[T any](val *T, options ...Option[T]) error {
	for i, o := range options {
		err := o(val)
		if err != nil {
			return fmt.Errorf("error in option %d: %w", i, err)
		}
	}
	return nil
}
