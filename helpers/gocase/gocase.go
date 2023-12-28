/*
 * The MIT License (MIT)
 *
 * Copyright (c) 2022 Takuo Oki
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, Subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or Substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

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

	"github.com/steebchen/prisma-client-go/helpers/strcase"
)

// ToLower returns a string converted to Go upper case.
func ToLower(s string) string {
	return defaultConverter.To(s, false)
}

// ToUpper returns a string converted to Go lower case.
func ToUpper(s string) string {
	return defaultConverter.To(s, true)
}

// To returns a string converted to Go case with converter.
func (c *Converter) To(s string, upper bool) string {
	for _, i := range c.initialisms {
		// not end
		re1 := regexp.MustCompile(fmt.Sprintf("%s([^a-z])", i.capUpper()))
		s = re1.ReplaceAllString(s, i.allUpper()+"$1")

		// end
		re2 := regexp.MustCompile(fmt.Sprintf("%s$", i.capUpper()))
		s = re2.ReplaceAllString(s, i.allUpper())
	}

	if len(s) == 0 {
		return s
	}

	var isAllCap = true
	for _, i := range s {
		if i < 'A' || i > 'Z' {
			isAllCap = false
		}
	}

	if strings.Contains(s, "_") || isAllCap {
		if upper {
			// for snake case
			s = strcase.ToUpperCamel(s)
		} else {
			s = strcase.ToLowerCamel(s)
		}
	} else {
		if upper {
			s = strings.ToUpper(s[:1]) + s[1:]
		} else {
			// TODO!!!
			s = strings.ToLower(s[:1]) + s[1:]
		}
	}

	// run again for new uppercase words
	for _, i := range c.initialisms {
		// not end
		re1 := regexp.MustCompile(fmt.Sprintf("%s([^a-z])", i.capUpper()))
		s = re1.ReplaceAllString(s, i.allUpper()+"$1")

		// end
		re2 := regexp.MustCompile(fmt.Sprintf("%s$", i.capUpper()))
		s = re2.ReplaceAllString(s, i.allUpper())
	}

	// fix casing after numbers
	n := strings.Builder{}
	n.Grow(len(s))
	prevIsNumber := false
	for _, v := range []byte(s) {
		vIsLow := v >= 'a' && v <= 'z'
		if prevIsNumber && vIsLow {
			v += 'A'
			v -= 'a'
		}
		prevIsNumber = v >= '0' && v <= '9'
		n.WriteByte(v)
	}

	return n.String()
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
