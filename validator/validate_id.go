package validator

import (
	"fmt"

	"github.com/knstch/subtrack-libs/svcerrs"
)

var (
	ErrBadID = fmt.Errorf("bad ID: %w", svcerrs.ErrInvalidData)
)

func ValidateID[T string | int | uint](id T) error {
	switch v := any(id).(type) {
	case int:
		if v < 1 {
			return ErrBadID
		}
	case uint:
		if v < 1 {
			return ErrBadID
		}
	case string:
		if v == "" || v == "0" {
			return ErrBadID
		}
	default:
		return ErrBadID
	}
	return nil
}
