#Pipeln

This is a negroni middleware inspired by Sprockets with Rails.

Curently it has the following features:
* Asset lookup based on a default (but configurable) layout.
* Processing Javascript //=require & //=require_tree
* Template helpers as convenience functions for script & link tags.

###Example usage
```go
package main

import (
        "net/http"

        "github.com/codegangsta/negroni"
        "github.com/losinggeneration/pipeln"
)

func main() {
        mux := http.NewServeMux()
        n := negroni.Classic()

        n.Use(pipeln.NewAssets())
        n.UseHandler(mux)

        n.Run(":8080")
}
```

###Middleware
The middleware  should be pretty straight forward. You tell negoroni to use pipeln.NewAssets and pass it an optional pipeln/assets.Options to change the directory layout if you want. The base layout NewAssets is expecting is:
```
- assets
  |- images
  |- javascripts
  |- stylesheets
```

###Assets
The assets directory can be searched for assets. It also provides the functionality to process Javscript files to support the
```
//= require <arg>
//= require_tree <dir>
```
directives.

###Templates
Calling FuncMap and passing an *html.Template as an argument will insert in several helpers to make some head tags for link & script easier. It uses the sprockets-rails names:
```
javascript_include_tag "script" tag_opts ( "attribute" "overrides" )
stylesheet_link_tag "stylesheet" tag_opts ( "attribute" "overrides" )
favicon_link_tag "icon" tag_opts ( "attribute" "overrides" )
tag_opts "attribute" "overrides"
```
tag_opts on the *_tag's are optional.


