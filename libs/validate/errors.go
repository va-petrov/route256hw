package validate

import (
	"github.com/pkg/errors"
)

var (
	ErrEmptyUser      = errors.New("empty user")
	ErrEmptyOrderID   = errors.New("empty orderID")
	ErrEmptySKU       = errors.New("empty sku")
	ErrEmptyItemsList = errors.New("empty order items list")
)

func Combine(errs ...error) error {
	var result error
	for _, err := range errs {
		if err != nil {
			if result != nil {
				result = errors.WithMessage(result, err.Error())
			} else {
				result = err
			}
		}
	}
	return result
}
