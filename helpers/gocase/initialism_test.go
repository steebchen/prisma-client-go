package gocase

import (
	"testing"
)

func TestCreateInitialisms(t *testing.T) {
	t.Parallel()

	cases := []struct {
		in      []string
		want    []initialism
		wantErr string
	}{
		{in: []string{"ACL"}, want: []initialism{newInitialism("ACL", "Acl")}},
		{in: []string{"api", "aSCII"}, want: []initialism{newInitialism("API", "Api"), newInitialism("ASCII", "Ascii")}},
		{in: []string{"UTF!"}, wantErr: "input \"UTF!\" is not alpha-numeric character"},
	}

	for _, c := range cases {
		r, err := createInitialisms(c.in...)
		if c.wantErr == "" {
			switch {
			case err != nil:
				t.Errorf("error must not be occurred: %v", err)
			case len(r) != len(c.want):
				t.Errorf("value length doesn't match: %d (want %d)", len(r), len(c.want))
			default:
				for i, w := range c.want {
					if r[i] != w {
						t.Errorf("value doesn't match: %v (want %v)", r[i], w)
					}
				}
			}
		} else {
			if err == nil {
				t.Error("error must be occurred")
			} else if err.Error() != c.wantErr {
				t.Errorf("error doesn't match: %v (want %s)", err, c.wantErr)
			}
		}
	}
}

func TestConvertToOnlyFirstLetterCapitalizedString(t *testing.T) {
	cases := []struct {
		in      string
		want    string
		wantErr string
	}{
		{in: "ACL", want: "Acl"},
		{in: "api", want: "Api"},
		{in: "aSCII", want: "Ascii"},
		{in: "cPu", want: "Cpu"},
		{in: "UTF8", want: "Utf8"},
		{in: "UTF!", wantErr: "input \"UTF!\" is not alpha-numeric character"},
		{in: "aa\xe2", wantErr: "input is not valid UTF-8"},
	}

	for _, c := range cases {
		r, err := convertToOnlyFirstLetterCapitalizedString(c.in)
		if c.wantErr == "" {
			if err != nil {
				t.Errorf("error must not be occurred: %v", err)
			} else if r != c.want {
				t.Errorf("value doesn't match: %s (want %s)", r, c.want)
			}
		} else {
			if err == nil {
				t.Error("error must be occurred")
			} else if err.Error() != c.wantErr {
				t.Errorf("error doesn't match: %v (want %s)", err, c.wantErr)
			}
		}
	}
}
