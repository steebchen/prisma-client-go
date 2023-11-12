// Package gocase is a package to convert normal CamelCase to Golang's CamelCase and vice versa.
// Golang's CamelCase means a string that takes into account to Go's common initialisms.
// For more details, please see [initialisms section] in [Staticcheck].
//
// [Staticcheck]: https://staticcheck.io/
// [initialisms section]: https://staticcheck.io/docs/configuration/options/#initialisms
package gocase

import (
	"fmt"
	"regexp"
	"strings"
)

// To returns a string converted to Go case.
func To(s string) string {
	return defaultConverter.To(s)
}

// To returns a string converted to Go case with converter.
func (c *Converter) To(s string) string {
	for _, i := range c.initialisms {
		// not end
		re1 := regexp.MustCompile(fmt.Sprintf("%s([^a-z])", i.capUpper()))
		s = re1.ReplaceAllString(s, i.allUpper()+"$1")

		// end
		re2 := regexp.MustCompile(fmt.Sprintf("%s$", i.capUpper()))
		s = re2.ReplaceAllString(s, i.allUpper())
	}
	return s
}

// Revert returns a string converted from Go case to normal case.
// Note that it is impossible to accurately determine the word break in a string of
// consecutive uppercase words, so the conversion maynot work as expected.
func Revert(s string) string {
	return defaultConverter.Revert(s)
}

// Revert returns a string converted from Go case to normal case with converter.
// Note that it is impossible to accurately determine the word break in a string of
// consecutive uppercase words, so the conversion maynot work as expected.
func (c *Converter) Revert(s string) string {
	for _, i := range c.initialisms {
		s = strings.ReplaceAll(s, i.allUpper(), i.capUpper())
	}
	return s
}
