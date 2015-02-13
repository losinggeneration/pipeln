package assets

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
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

	// Store matches here
	matches := make([]map[string]string, 0)

	// file's contents
	var fsrc string

	// String to return with parsed files
	var src string

	for {
		// Read each line
		line, _, err := fbio.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return src, err
		}

		fsrc += string(line) + "\n"

		// Match ^//= require(_tree)? (.*)$
		r, err := regexp.Compile(`^//=[[:blank:]]*require(?P<tree>_tree)?[[:blank:]]+(?P<argument>.*)$`)
		if err != nil {
			return src, err
		}

		sm := r.FindStringSubmatch(string(line))
		m := make(map[string]string)
		if len(sm) > 0 {
			// Take each match and associate it with the named match
			for i, name := range r.SubexpNames() {
				// Don't include non-matches
				if sm[i] != "" {
					m[name] = sm[i]
				}
			}

			// Store the match from this line
			matches = append(matches, m)
		}
	}

	// Go through each match and process it appropriately
	for _, m := range matches {
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
	}

	src += fsrc

	return src, nil
}

func javascriptRequireTree(filename, argument string) (string, error) {
	return "", nil
}

func javascriptRequire(filename, argument string) (string, error) {
	file, err := os.Open(filepath.Join(filepath.Dir(filename), argument))
	if err != nil {
		return "", err
	}
	defer file.Close()

	src, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(src), nil
}
