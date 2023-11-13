package gocase

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	runeOffsetUpperToLower = 32
)

// initialism is a type that describes initialization rule.
// The first element is set to an all uppercase string.
// The second element is set to a string with only the first letter capitalized.
type initialism [2]string

func newInitialism(s1, s2 string) initialism {
	return [2]string{s1, s2}
}

func (i initialism) allUpper() string {
	return i[0]
}

func (i initialism) capUpper() string {
	return i[1]
}

func createInitialisms(initialisms ...string) ([]initialism, error) {
	results := make([]initialism, 0, len(initialisms))
	for _, i := range initialisms {

		s, err := convertToOnlyFirstLetterCapitalizedString(i)
		if err != nil {
			return nil, err
		}

		results = append(results, newInitialism(strings.ToUpper(i), s))
	}
	return results, nil
}

func convertToOnlyFirstLetterCapitalizedString(str string) (string, error) {
	if !utf8.ValidString(str) {
		return "", errors.New("input is not valid UTF-8")
	}

	var result []rune
	for i, r := range str {
		switch {
		case 'A' <= r && r <= 'Z':
			if i == 0 {
				result = append(result, r)
			} else {
				result = append(result, rune(int(r)+runeOffsetUpperToLower))
			}
		case 'a' <= r && r <= 'z':
			if i == 0 {
				result = append(result, rune(int(r)-runeOffsetUpperToLower))
			} else {
				result = append(result, r)
			}
		case '0' <= r && r <= '9':
			result = append(result, r)
		default:
			return "", fmt.Errorf("input %q is not alpha-numeric character", str)
		}
	}

	return string(result), nil
}
