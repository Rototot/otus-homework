package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")
var errUnexpectedValue = errors.New("unexpected value")

const (
	// "\"
	escapeSymbol rune = 92
)

func Unpack(rawData string) (string, error) {
	var rawDataRunes = []rune(rawData)
	var unpackedData strings.Builder

	for i := 0; i < len(rawDataRunes); {
		var err error
		var current = rawDataRunes[i]
		var next = nextRune(rawDataRunes, i+1)

		// digits without escape - cannot be processed
		if unicode.IsDigit(current) {
			return "", ErrInvalidString
		}

		// decode escaping
		if current == escapeSymbol {
			//decode \rune -> rune
			current, err = decode(next)
			if err != nil {
				return "", err
			}

			next = nextRune(rawDataRunes, i+2) //nolint:gomnd
			i++
		}

		// repeat runes
		// zero qty - write as allowed symbol
		if unicode.IsDigit(next) {
			qty, err := getRepeatableQty(next)
			if err != nil {
				return "", err
			}

			_, err = unpackedData.WriteString(strings.Repeat(string(current), qty))
			if err != nil {
				return "", err
			}
			// skip repeatable value
			i++
		} else {
			// write allowed symbol
			unpackedData.WriteRune(current)
		}

		i++
	}

	return unpackedData.String(), nil
}

func getRepeatableQty(value rune) (int, error) {
	if !unicode.IsDigit(value) {
		return 0, errUnexpectedValue
	}

	qty, err := strconv.Atoi(string(value))
	if err != nil {
		return 0, err
	}

	// write 1 symbol anywhere
	if qty == 0 {
		return 1, nil
	}

	return qty, nil
}

func nextRune(data []rune, index int) rune {
	if len(data) > index {
		return data[index]
	}

	return 0
}

func decode(value rune) (rune, error) {
	// escape only "\"|[0-9]
	if unicode.IsDigit(value) || value == escapeSymbol {
		return value, nil
	}

	return 0, ErrInvalidString
}
