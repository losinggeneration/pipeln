package template

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"strconv"
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

func verifyWant(got, want string) error {
	if got == want {
		return nil
	}

	// Make sure it's not missing any expected tags
	for _, s := range strings.Split(want, " ") {
		if !strings.Contains(got, s) {
			return errors.New(fmt.Sprintf("'%v' is not in '%v'", s, got))
		}
	}

	return nil
}

func verifyGot(got, want string) error {
	if got == want {
		return nil
	}

	// Make sure all tags we got are there
	for _, s := range strings.Split(got, " ") {
		if !strings.Contains(want, s) {
			return errors.New(fmt.Sprintf("'%v' is not in '%v'", s, want))
		}
	}

	return nil
}

func verifyTag(got, want string) error {
	if err := verifyGot(got, want); err != nil {
		return err
	}

	if err := verifyWant(got, want); err != nil {
		return err
	}

	return nil
}

func TestFuncMap(t *testing.T) {
	tests := []struct {
		template string
		want     string
	}{
		{`{{javascript_include_tag "a"}}`, `<script src="a" ></script>`},
		{`{{javascript_include_tag "a" ( tag_opts "c" "d" )}}`, `<script src="a" c="d" ></script>`},
		{`{{stylesheet_link_tag "a"}}`, `<link href="a" rel="stylesheet" >`},
		{`{{stylesheet_link_tag "a" ( tag_opts "c" "d" )}}`, `<link href="a" rel="stylesheet" c="d" >`},
		{`{{favicon_link_tag "a"}}`, `<link href="a" type="image/vnd.microsoft.icon" >`},
		{`{{favicon_link_tag "a" ( tag_opts "type" "image/png" "rel" "shortcut" )}}`, `<link href="a" type="image/png" rel="shortcut" >`},
	}

	for i, s := range tests {
		temp := template.New("test" + strconv.Itoa(i))
		FuncMap(temp)

		temp, err := temp.Parse(s.template)
		if err != nil {
			fmt.Println(s.template)
			t.Error(err)
		}

		buffer := bytes.NewBuffer(make([]byte, 0))
		if err := temp.Execute(buffer, nil); err != nil {
			t.Error(err)
		}

		got := buffer.String()
		if err := verifyWant(got, s.want); err != nil {
			t.Error(err)
		}
	}
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

func TestJavascriptIncludeTag(t *testing.T) {
	tests := []struct {
		tag  string
		args []string
		want string
	}{
		{"a.js", nil, `<script src="a.js" ></script>`},
		{"a.js", []string{"fu", "bar"}, `<script src="a.js" fu="bar" ></script>`},
		{"a.js", []string{"fu", "bar", "bar", "fu"}, `<script src="a.js" fu="bar" bar="fu" ></script>`},
	}

	for _, s := range tests {
		tagopt, err := newTagOpts(s.args...)
		if err != nil {
			t.Error(err)
		}

		got := string(JavascriptIncludeTag(s.tag, tagopt))
		if err := verifyTagOpts(got, tagopt); err != nil {
			t.Error(err)
		}

		if err := verifyWant(got, s.want); err != nil {
			t.Error(err)
		}
	}
}

func TestLinkTag(t *testing.T) {
	tests := []struct {
		tag  string
		args []string
		want string
	}{
		{"a", nil, `<link href="a" >`},
		{"b", []string{"fu", "bar"}, `<link href="b" fu="bar" >`},
		{"c", []string{"fu", "bar", "bar", "fu"}, `<link href="c" fu="bar" bar="fu" >`},
	}

	for _, s := range tests {
		tagopt, err := newTagOpts(s.args...)
		if err != nil {
			t.Error(err)
		}

		got := string(linkTag(s.tag, tagopt))
		if err := verifyTagOpts(got, tagopt); err != nil {
			t.Error(err)
		}

		if err := verifyWant(got, s.want); err != nil {
			t.Error(err)
		}
	}
}

func TestStylesheetLinkTag(t *testing.T) {
	tests := []struct {
		tag  string
		args []string
		want string
	}{
		{"a", nil, `<link href="a" rel="stylesheet" >`},
		{"b", []string{"fu", "bar"}, `<link href="b" rel="stylesheet" fu="bar" >`},
		{"c", []string{"fu", "bar", "bar", "fu"}, `<link href="c" rel="stylesheet" fu="bar" bar="fu" >`},
	}

	for _, s := range tests {
		tagopt, err := newTagOpts(s.args...)
		if err != nil {
			t.Error(err)
		}

		got := string(StylesheetLinkTag(s.tag, tagopt))
		if err := verifyTagOpts(got, tagopt); err != nil {
			t.Error(err)
		}

		if err := verifyWant(got, s.want); err != nil {
			t.Error(err)
		}
	}
}

func TestFaviconLinkTag(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"", nil, `<link href="" type="image/vnd.microsoft.icon" >`},
		{"icon.ico", nil, `<link href="icon.ico" type="image/vnd.microsoft.icon" >`},
		{"icon.jpg", []string{"rel", "shortcut"}, `<link href="icon.jpg" type="image/vnd.microsoft.icon" rel="shortcut" >`},
		{"icon.png", []string{"rel", "shortcut", "type", "image/png"}, `<link href="icon.png" type="image/png" rel="shortcut" >`},
	}

	for _, s := range tests {
		tagopt, err := newTagOpts(s.args...)
		if err != nil {
			t.Error(err)
		}

		got := string(FaviconLinkTag(s.name, tagopt))
		if err := verifyTagOpts(got, tagopt); err != nil {
			t.Error(err)
		}

		if err := verifyWant(got, s.want); err != nil {
			t.Error(err)
		}
	}
}
