package gocase_test

import (
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/steebchen/prisma-client-go/helpers/gocase"
)

func TestConverter_ToLower(t *testing.T) {
	t.Parallel()

	dc, _ := gocase.New()
	cc, _ := gocase.New(gocase.WithInitialisms("JSON", "CSV"))

	cases := []struct {
		conv       *gocase.Converter
		have, want string
	}{
		{conv: dc, have: "", want: ""},
		{conv: dc, have: "CONSTANT", want: "constant"},
		{conv: dc, have: "id", want: "id"},
		{conv: dc, have: "ID", want: "id"},
		{conv: dc, have: "jsonFile", want: "jsonFile"},
		// {conv: dc, have: "IpAddress", want: "ipAddress"},
		{conv: dc, have: "ip_address", want: "ipAddress"},
		{conv: dc, have: "defaultDnsServer", want: "defaultDNSServer"},
		{conv: dc, have: "somethingHttpApiId", want: "somethingHTTPAPIID"},
		{conv: dc, have: "somethingUuid", want: "somethingUUID"},
		{conv: dc, have: "somethingSip", want: "somethingSIP"},
		{conv: dc, have: "Urid", want: "urid"},
		{conv: dc, have: "stuffLast7D", want: "stuffLast7D"},
		{conv: dc, have: "stuffLast7d", want: "stuffLast7D"},
		{conv: dc, have: "StuffLast7d", want: "stuffLast7D"},
		{conv: dc, have: "StuffLast7dAnd", want: "stuffLast7DAnd"},
		{conv: dc, have: "StuffLast7DAnd", want: "stuffLast7DAnd"},
		{conv: dc, have: "anotherIDStuffSomethingID", want: "anotherIDStuffSomethingID"},
		{conv: dc, have: "anotherIdStuffSomethingId", want: "anotherIDStuffSomethingID"},
		{conv: dc, have: "anotherIdStuffSomethingId", want: "anotherIDStuffSomethingID"},
		{conv: dc, have: "another_id_stuff_something_id", want: "anotherIDStuffSomethingID"},
		// {conv: dc, have: "APISession", want: "apiSession"},

		// {conv: cc, have: "JsonFile", want: "jsonFile"},
		// {conv: cc, have: "CsvFile", want: "csvFile"},
		{conv: cc, have: "IpAddress", want: "ipAddress"},
	}

	for _, c := range cases {
		cc := c
		t.Run(fmt.Sprintf("%s -> %s", cc.have, cc.want), func(t *testing.T) {
			r := cc.conv.To(cc.have, false)
			if r != cc.want {
				t.Errorf("value doesn't match: have %s, is %s, want %s", cc.have, r, cc.want)
			}
		})
	}
}

func TestConverter_ToUpper(t *testing.T) {
	t.Parallel()

	dc, _ := gocase.New()
	cc, _ := gocase.New(gocase.WithInitialisms("JSON", "CSV"))

	cases := []struct {
		conv       *gocase.Converter
		have, want string
	}{
		{conv: dc, have: "", want: ""},
		{conv: dc, have: "CONSTANT", want: "Constant"},
		{conv: dc, have: "id", want: "ID"},
		{conv: dc, have: "Id", want: "ID"},
		{conv: dc, have: "IdSomething", want: "IDSomething"},
		{conv: dc, have: "IDSomething", want: "IDSomething"},
		{conv: dc, have: "jsonFile", want: "JSONFile"},
		{conv: dc, have: "IpAddress", want: "IPAddress"},
		{conv: dc, have: "ip_address", want: "IPAddress"},
		{conv: dc, have: "defaultDnsServer", want: "DefaultDNSServer"},
		{conv: dc, have: "somethingHttpApiId", want: "SomethingHTTPAPIID"},
		{conv: dc, have: "somethingUuid", want: "SomethingUUID"},
		{conv: dc, have: "somethingSip", want: "SomethingSIP"},
		{conv: dc, have: "stuffLast7D", want: "StuffLast7D"},
		{conv: dc, have: "stuffLast7d", want: "StuffLast7D"},
		{conv: dc, have: "StuffLast7d", want: "StuffLast7D"},
		{conv: dc, have: "StuffLast7dAnd", want: "StuffLast7DAnd"},
		{conv: dc, have: "StuffLast7DAnd", want: "StuffLast7DAnd"},
		{conv: dc, have: "Urid", want: "Urid"},
		{conv: dc, have: "anotherIDStuffSomethingID", want: "AnotherIDStuffSomethingID"},
		{conv: dc, have: "anotherIdStuffSomethingId", want: "AnotherIDStuffSomethingID"},
		{conv: dc, have: "anotherIdStuffSomethingId", want: "AnotherIDStuffSomethingID"},
		{conv: dc, have: "another_id_stuff_something_id", want: "AnotherIDStuffSomethingID"},
		{conv: dc, have: "APISession", want: "APISession"},

		{conv: cc, have: "JsonFile", want: "JSONFile"},
		{conv: cc, have: "CsvFile", want: "CSVFile"},
		{conv: cc, have: "IpAddress", want: "IpAddress"},
	}

	for _, c := range cases {
		cc := c
		t.Run(fmt.Sprintf("%s -> %s", cc.have, cc.want), func(t *testing.T) {
			r := cc.conv.To(cc.have, true)
			if r != cc.want {
				t.Errorf("value doesn't match: have %s, is %s, want %s", cc.have, r, cc.want)
			}
		})
	}
}

func TestConverter_Revert(t *testing.T) {
	dc, _ := gocase.New()
	cc, _ := gocase.New(gocase.WithInitialisms("JSON", "CSV"))

	cases := []struct {
		conv    *gocase.Converter
		s, want string
	}{
		{conv: dc, s: "", want: ""},
		{conv: dc, s: "jsonFile", want: "jsonFile"},
		{conv: dc, s: "IPAddress", want: "IpAddress"},
		{conv: dc, s: "defaultDNSServer", want: "defaultDnsServer"},
		{conv: dc, s: "somethingHTTPAPIID", want: "somethingHttpApiId"},
		{conv: dc, s: "somethingUUID", want: "somethingUuid"},
		{conv: dc, s: "somethingSIP", want: "somethingSip"},
		{conv: cc, s: "JSONFile", want: "JsonFile"},
		{conv: cc, s: "CSVFile", want: "CsvFile"},
		{conv: cc, s: "somethingSIP", want: "somethingSIP"},
	}

	for _, c := range cases {
		r := c.conv.Revert(c.s)
		if r != c.want {
			t.Errorf("value doesn't match: %s (want %s)", r, c.want)
		}
	}
}

// FuzzReverseUpper runs a Fuzzing test to check if the strings
// before and after `To` and `Revert` match.
// Note that there may be cases where the strings before and after
// the `To` and `Revert` do not match for certain inputs.
//
// ```cmd
// go test -fuzz=Fuzz
// ```
func FuzzReverseUpper(f *testing.F) {
	testcases := []string{"JsonFile", "IpAddress", "DefaultDnsServer"}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, orig string) {
		to := gocase.ToUpper(orig)
		rev := gocase.Revert(to)
		if !ignoreInput(orig) && orig != rev {
			t.Errorf("before: %q, after: %q", orig, rev)
		}
		if utf8.ValidString(orig) && !utf8.ValidString(rev) {
			t.Errorf("To or Revert produced invalid UTF-8 string %q", rev)
		}
	})
}

// FuzzReverseLower runs a Fuzzing test to check if the strings
// before and after `To` and `Revert` match.
// Note that there may be cases where the strings before and after
// the `To` and `Revert` do not match for certain inputs.
func FuzzReverseLower(f *testing.F) {
	testcases := []string{"jsonFile", "ipAddress", "defaultDnsServer"}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, orig string) {
		to := gocase.ToLower(orig)
		rev := gocase.Revert(to)
		if !ignoreInput(orig) && orig != rev {
			t.Errorf("before: %q, after: %q", orig, rev)
		}
		if utf8.ValidString(orig) && !utf8.ValidString(rev) {
			t.Errorf("To or Revert produced invalid UTF-8 string %q", rev)
		}
	})
}

func ignoreInput(in string) bool {

	for _, s := range gocase.DefaultInitialisms {
		if strings.Contains(in, s) {
			return true
		}
	}

	return false
}
