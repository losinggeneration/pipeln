package template

import (
	"errors"
	"html/template"
	"strings"
)

var FuncMap = template.FuncMap{
	"javascript_include_tag": JavascriptIncludeTag,
	"stylesheet_link_tag":    StylesheetLinkTag,
	"favicon_link_tag":       FaviconLinkTag,
	"tag_opts":               newTagOpts,
}

/*
 * Funcs adds functions to the template
 *
 * javascript_include_tag
 * Adds a script tag.
 * Arguments:
 * name
 * ( tag_opts )
 *
 * stylesheet_link_tag
 * Adds a script tag.
 * Arguments:
 * name
 * ( tag_opts )
 *
 * favicon_link_tag
 * Adds a link rel="icon" tag.
 * Arguments:
 * name
 * ( tag_opts )
 *
 * tag_opts
 * Adds attributes to a tag.
 * Arguments:
 * space separated tuples
 * e.g. "rel" "icon" "type" "image/png"
 * Returns an error stopping the template parsing if an odd number of
 * arguments are given
 *
 * It's important to use parens separated by spaces around the tag_opts
 * when used to pass arguments to the various tags.
 * e.g. javascript_include_tag "bar" ( tag_opts "charset" "utf-8" )
 */
func Funcs(t *template.Template) {
	t.Funcs(FuncMap)
}

/* TagOpts are attribute options to insert into a tag */
type TagOpts map[string]string

/* String converts the TagOpts map to an attribute string */
func (t TagOpts) String() string {
	var to string

	for k, v := range t {
		to += k + `="` + v + `" `
	}

	return to
}

/* newTagOpts creates tag attributes */
func newTagOpts(opts ...string) (TagOpts, error) {
	var to TagOpts

	if len(opts) == 0 {
		return make(TagOpts), nil
	} else if len(opts)%2 == 1 {
		return to, errors.New("expects an even number of parameters")
	} else {
		to = make(TagOpts)
	}

	for i := 0; i < len(opts); i += 2 {
		to[opts[i]] = opts[i+1]
	}

	return to, nil
}

func getTagOpt(opts ...TagOpts) TagOpts {
	var o TagOpts
	if len(opts) == 0 {
		o = make(TagOpts)
	} else {
		o = opts[0]
	}

	return o
}

/* tag creates an open tag */
func tag(t string, opts ...TagOpts) template.HTML {
	tag := "<" + t + " "
	if len(opts) > 0 {
		tag += opts[0].String()
	}
	tag = strings.TrimSpace(tag)
	tag += ">"

	return template.HTML(tag)
}

/* JavascriptIncludeTag creates a script src tag */
func JavascriptIncludeTag(name string, opts ...TagOpts) template.HTML {
	o := getTagOpt(opts...)

	o["src"] = name

	return tag("script", o) + "</script>"
}

func linkTag(name string, opts ...TagOpts) template.HTML {
	o := getTagOpt(opts...)
	o["href"] = name

	return tag("link", o)
}

/* StylesheetLinkTag creates a link rel="stylesheet" tag */
func StylesheetLinkTag(name string, opts ...TagOpts) template.HTML {
	o := getTagOpt(opts...)

	o["rel"] = "stylesheet"

	return linkTag(name, o)
}

/* FaviconLinkTag creates a link rel="shortcut icon" tag */
func FaviconLinkTag(name string, opts ...TagOpts) template.HTML {
	o := getTagOpt(opts...)

	if o["rel"] == "" {
		o["rel"] = "shortcut icon"
	}
	if o["type"] == "" {
		o["type"] = "image/vnd.microsoft.icon"
	}

	return linkTag(name, o)
}
