package engine

import (
	"errors"
	"strconv"
)

// ErrInvalid is an input that breaks a rule the engine relies on (docs/architecture.md, section
// 3.6): a choice is never made of it. It comes as a value and never by panicking; errors.Is tells
// it apart, and the message says what and where. An error of learner-model the input brought, from
// Readout.Task, is wrapped, so that errors.Is finds learner-model's kind too.
var ErrInvalid = errors.New("engine: invalid input")

// detailError is ErrInvalid with what and where, and the error that caused it, if any.
type detailError struct {
	detail string
	cause  error
}

func (e *detailError) Error() string {
	if e.cause != nil {
		return ErrInvalid.Error() + ": " + e.detail + ": " + e.cause.Error()
	}
	return ErrInvalid.Error() + ": " + e.detail
}

func (e *detailError) Unwrap() []error {
	if e.cause != nil {
		return []error{ErrInvalid, e.cause}
	}
	return []error{ErrInvalid}
}

func invalid(detail string) error { return &detailError{detail: detail} }

// shown quotes an ID for a message, byte by byte: printable ASCII as it is, every other byte as
// \x and two hex digits, so that the message reads the same whatever the bytes and whatever Go's
// Unicode tables.
func shown(id string) string {
	const hex = "0123456789abcdef"
	b := make([]byte, 0, len(id)+2)
	b = append(b, '"')
	for i := range len(id) {
		c := id[i]
		switch {
		case c == '"' || c == '\\':
			b = append(b, '\\', c)
		case c >= 0x20 && c < 0x7f:
			b = append(b, c)
		default:
			b = append(b, '\\', 'x', hex[c>>4], hex[c&0xf])
		}
	}
	return string(append(b, '"'))
}

func itoa(v int64) string { return strconv.FormatInt(v, 10) }
