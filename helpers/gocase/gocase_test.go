package gocase_test

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/steebchen/prisma-client-go/helpers/gocase"
)

func TestConverter_To(t *testing.T) {
	t.Parallel()

	dc, _ := gocase.New()
	cc, _ := gocase.New(gocase.WithInitialisms("JSON", "CSV"))

	cases := []struct {
		conv    *gocase.Converter
		s, want string
	}{
		{conv: dc, s: "", want: ""},
		{conv: dc, s: "jsonFile", want: "jsonFile"},
		{conv: dc, s: "IpAddress", want: "IPAddress"},
		{conv: dc, s: "defaultDnsServer", want: "defaultDNSServer"},
		{conv: dc, s: "somethingHttpApiId", want: "somethingHTTPAPIID"},
		{conv: dc, s: "somethingUuid", want: "somethingUUID"},
		{conv: dc, s: "somethingSip", want: "somethingSIP"},
		{conv: dc, s: "Urid", want: "Urid"},
		{conv: cc, s: "JsonFile", want: "JSONFile"},
		{conv: cc, s: "CsvFile", want: "CSVFile"},
		{conv: cc, s: "IpAddress", want: "IpAddress"},
	}

	for _, c := range cases {
		r := c.conv.To(c.s)
		if r != c.want {
			t.Errorf("value doesn't match: %s (want %s)", r, c.want)
		}
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

// FuzzReverse runs a Fuzzing test to check if the strings
// before and after `To` and `Revert` match.
// Note that there may be cases where the strings before and after
// the `To` and `Revert` do not match for certain inputs.
//
// ```cmd
// go test -fuzz=Fuzz
// ```
func FuzzReverse(f *testing.F) {
	testcases := []string{"jsonFile", "IpAddress", "defaultDnsServer"}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, orig string) {
		to := gocase.To(orig)
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
