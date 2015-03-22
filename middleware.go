package pipeln

import (
	"html/template"
	"io"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/losinggeneration/pipeln/assets"
	helper "github.com/losinggeneration/pipeln/html/template"
)

var FuncMap = helper.FuncMap

type Assets struct {
	assets.Assets
}

func javascriptIncludeTag(name string, opts ...helper.TagOpts) template.HTML {
	return helper.JavascriptIncludeTag(name, opts...)
}

func stylesheetLinkTag(name string, opts ...helper.TagOpts) template.HTML {
	return helper.StylesheetLinkTag(name, opts...)
}

func faviconLinkTag(name string, opts ...helper.TagOpts) template.HTML {
	return helper.FaviconLinkTag(name, opts...)
}

func Funcs(t *template.Template) {
	helper.Funcs(t)

	// override these with the ones that lookup the assets
	t.Funcs(template.FuncMap{
		"javascript_include_tag": javascriptIncludeTag,
		"stylesheet_link_tag":    stylesheetLinkTag,
		"favicon_link_tag":       faviconLinkTag,
	})
}

func NewAssets(options ...assets.Options) *Assets {
	var o assets.Options
	if len(options) == 0 {
		o = assets.Options{
			BaseDir:        "assets",
			JavascriptsDir: "javascripts",
			StylesheetsDir: "stylesheets",
			ImagesDir:      "imgages",
		}
	} else {
		o = options[0]
	}

	return &Assets{
		assets.Assets{Options: o},
	}
}

func (a *Assets) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// Don't process if nothing is passed in
	if req.RequestURI == "/" {
		next(w, req)
		return
	}

	file, err := a.Lookup(req.RequestURI)
	if err != nil {
		if err == assets.ErrNotFound {
			next(w, req)
			return
		} else {
			log.Println(err)
			return
		}
	}
	defer file.Close()

	processed, err := a.Process(req.RequestURI)
	if err == nil && processed != "" {
		w.Header().Add("Content-Type", mime.TypeByExtension(filepath.Ext(file.Name())))
		w.Header().Add("Content-Length", strconv.Itoa(len(processed)))
		io.Copy(w, strings.NewReader(processed))
		return
	}

	fi, err := file.Stat()
	if err != nil {
		log.Println(err)
		next(w, req)
		return
	}

	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}
