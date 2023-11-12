package gocase

import (
	"testing"
)

func TestTo(t *testing.T) {
	cases := []struct {
		s, expected string
	}{
		{s: "", expected: ""},
		{s: "jsonFile", expected: "jsonFile"},
		{s: "IpAddress", expected: "IPAddress"},
		{s: "defaultDnsServer", expected: "defaultDNSServer"},
		{s: "somethingHttpApiId", expected: "somethingHTTPAPIID"},
	}

	for _, c := range cases {
		r := To(c.s)
		if r != c.expected {
			t.Errorf("value doesn't match: %s (expected %s)", r, c.expected)
		}
	}
}

func TestRevert(t *testing.T) {
	cases := []struct {
		s, expected string
	}{
		{s: "", expected: ""},
		{s: "jsonFile", expected: "jsonFile"},
		{s: "IPAddress", expected: "IpAddress"},
		{s: "defaultDNSServer", expected: "defaultDnsServer"},
		{s: "somethingHTTPAPIID", expected: "somethingHttpApiId"},
	}

	for _, c := range cases {
		r := Revert(c.s)
		if r != c.expected {
			t.Errorf("value doesn't match: %s (expected %s)", r, c.expected)
		}
	}
}
