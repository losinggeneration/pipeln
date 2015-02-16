package template

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func verifyTagOpts(want string, opts ...TagOpts) error {
	if len(opts) == 0 || opts == nil {
		return nil
	}

	for _, opt := range opts {
		for k, v := range opt {
			if strings.Contains(want, fmt.Sprintf("%v=%q", k, v)) == false {
				return errors.New(fmt.Sprintf("%v=%q is not in '%v'", k, v, want))
			}
		}
	}

	return nil
}

func verifyTag(got, want string) error {
	if got == want {
		return nil
	}

	// Make sure all tags we got are there
	for _, s := range strings.Split(got, " ") {
		if !strings.Contains(want, s) {
			return errors.New(fmt.Sprintf("'%v' is not in '%v'", s, want))
		}
	}

	// Make sure it's not missing any expected tags
	for _, s := range strings.Split(want, " ") {
		if !strings.Contains(got, s) {
			return errors.New(fmt.Sprintf("'%v' is not in '%v'", s, got))
		}
	}
	return nil
}

func TestTagOptsStringer(t *testing.T) {
	tests := []struct {
		tagopt TagOpts
		want   string
	}{
		{TagOpts{"fu": "bar"}, `fu="bar" `},
		{TagOpts{"fu": "bar", "bar": "fu"}, `fu="bar" bar="fu" `},
		{TagOpts{}, ``},
	}

	for _, s := range tests {
		if err := verifyTagOpts(s.want, s.tagopt); err != nil {
			t.Error(err)
		}
	}
}

func TestNewTagOpts(t *testing.T) {
	tests := []struct {
		args []string
		want string
	}{
		{nil, ""},
		{[]string{}, ""},
		{[]string{"fu", "bar"}, `fu="bar" `},
		{[]string{"fu", "bar", "bar", "fu"}, `fu="bar" bar="fu" `},
	}

	for _, s := range tests {
		tagopt, err := newTagOpts(s.args...)
		if err != nil {
			t.Error(err)
		}

		if err := verifyTagOpts(tagopt.String(), tagopt); err != nil {
			t.Error(err)
		}

		if err := verifyTagOpts(s.want, tagopt); err != nil {
			t.Error(err)
		}
	}

	errors := [][]string{
		[]string{"", "", ""},
	}

	for _, s := range errors {
		tagopt, err := newTagOpts(s...)
		if err == nil {
			t.Errorf("Expected an error with: %q. Got: '%v'", s, tagopt.String())
		}
	}
}

func TestGetTagOpt(t *testing.T) {
	t1 := make(TagOpts)
	t2, _ := newTagOpts("fu", "bar")
	t3, _ := newTagOpts("fu", "bar", "bar", "fu")

	tests := []struct {
		tagopts []TagOpts
		want    string
	}{
		{[]TagOpts{}, ""},
		{[]TagOpts{nil, t2}, ""},
		{[]TagOpts{t1, t2}, ""},
		{[]TagOpts{t2, t1}, t2.String()},
		{[]TagOpts{t2, t3, t1, t2, t1}, t2.String()},
	}

	for _, s := range tests {
		tagopt := getTagOpt(s.tagopts...)

		if err := verifyTagOpts(tagopt.String(), tagopt); err != nil {
			t.Error(err)
		}

		if err := verifyTagOpts(s.want, tagopt); err != nil {
			t.Error(err)
		}
	}
}

func TestTag(t *testing.T) {
	tests := []struct {
		tag  string
		args []string
		want string
	}{
		{"div", nil, "<div>"},
		{"div", []string{"fu", "bar"}, `<div fu="bar">`},
		{"div", []string{"fu", "bar", "bar", "fu"}, `<div fu="bar" bar="fu">`},
	}

	for _, s := range tests {
		tagopt, err := newTagOpts(s.args...)
		if err != nil {
			t.Error(err)
		}

		got := string(tag(s.tag, tagopt))
		if err := verifyTagOpts(got, tagopt); err != nil {
			t.Error(err)
		}

		if err := verifyTag(got, s.want); err != nil {
			t.Error(err)
		}
	}
}
