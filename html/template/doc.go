/*
This provides some helpers for using within templates. It allows you to
quite easily include scripts & stylesheets within a template.

Example

import "html/template"
import helpers "pipeln/html/template"
...
t, err := template.New("foo").Parse(`<html><head>
{{stylesheet_link_tag "foo.css"}}
{{javascript_include_tag "foo.js"}}
{{favicon_link_tag "logo.png" ( tag_opts "type" "image/png" )}}
</head></html>`)
err = t.ExecuteTemplate(out, nil)

Produces
<html><head>
<script src="foo.js"></script>
<link rel="stylesheet" href="foo.js">
<link rel="shortcut icon" type="image/png" href="logo.png">
</head></html>
*/
package template
