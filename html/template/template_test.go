package template

import "testing"

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
		if s.tagopt.String() != s.want {
			t.Errorf(`'%v' != '%v'`, s.tagopt.String(), s.want)
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

		if tagopt.String() != s.want {
			t.Errorf(`'%v' != '%v'`, tagopt.String(), s.want)
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
		if tagopt.String() != s.want {
			t.Errorf(`'%v' != '%v'`, tagopt.String(), s.want)
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

		if string(tag(s.tag, tagopt)) != s.want {
			t.Errorf(`'%v' != '%v'`, tagopt.String(), s.want)
		}
	}
}
