package assets

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Options struct {
	BaseDir        string
	JavascriptsDir string
	StylesheetsDir string
	ImagesDir      string
}

type Assets struct {
	Options Options
}

func (a Assets) Lookup(name string) (*os.File, error) {
	return Lookup(name, a.Options)
}

func (a Assets) Process(name string) (string, error) {
	return Process(name, a.Options)
}

func (a Assets) ProcessFile(f *os.File) (string, error) {
	return ProcessFile(f)
}

var options = Options{
	BaseDir:        "assets",
	JavascriptsDir: "javascripts",
	StylesheetsDir: "stylesheets",
	ImagesDir:      "images",
}

var (
	ErrNotFound = errors.New("asset not found")
)

func Lookup(name string, opts ...Options) (*os.File, error) {
	if opts == nil {
		opts = make([]Options, 1)
		opts[0] = options
	}

	path := filepath.Join(opts[0].BaseDir, name)
	file, err := os.Open(path)

	if err != nil {
		return nil, ErrNotFound
	}

	return file, nil
}

func Process(name string, opts ...Options) (string, error) {
	f, err := Lookup(name, opts...)
	if err != nil {
		return "", err
	}

	return ProcessFile(f)
}

func ProcessFile(f *os.File) (string, error) {
	defer f.Close()

	switch {
	case strings.Contains(f.Name(), ".js"):
		return processJavascriptRequires(f)
	}

	return "", nil
}

func processJavascriptRequires(f *os.File) (string, error) {
	fbio := bufio.NewReader(f)

	// file's contents
	var fsrc string

	for {
		// Read each line
		line, _, err := fbio.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return fsrc, err
		}

		fsrc += string(line) + "\n"

		// Match ^//= require(_tree)? (.*)$
		r, err := regexp.Compile(`^//=[[:blank:]]*require(?P<tree>_tree)?[[:blank:]]+(?P<argument>.*)$`)
		if err != nil {
			return fsrc, err
		}

		sm := r.FindStringSubmatch(string(line))
		m := make(map[string]string)
		if len(sm) > 0 {
			var src string
			// Take each match and associate it with the named match
			for i, name := range r.SubexpNames() {
				// Don't include non-matches
				if sm[i] != "" {
					m[name] = sm[i]
				}
			}
			if m["tree"] != "" {
				s, err := javascriptRequireTree(f.Name(), m["argument"])
				if err != nil {
					return src, err
				}
				src += s
			} else {
				s, err := javascriptRequire(f.Name(), m["argument"])
				if err != nil {
					return src, err
				}
				src += s
			}
			fsrc += src
		}
	}

	return fsrc, nil
}

func javascriptRequireTree(filename, argument string) (string, error) {
	var src string
	base := filepath.Join(filepath.Dir(filename), argument)

	filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		// don't process directories as if they're files
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".js" {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			s, err := processJavascriptRequires(file)
			if err != nil {
				return err
			}
			src += s
		}

		return nil
	})

	return src, nil
}

func javascriptRequire(filename, argument string) (string, error) {
	if !strings.Contains(argument, ".js") {
		argument += ".js"
	}

	file, err := os.Open(filepath.Join(filepath.Dir(filename), argument))
	if err != nil {
		return "", err
	}
	defer file.Close()

	return processJavascriptRequires(file)
}
