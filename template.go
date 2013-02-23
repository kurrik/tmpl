// Copyright 2013 Arne Roomann-Kurrik
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tmpl

import (
	"io"
	"bytes"
	"fmt"
	"github.com/kurrik/fauxfile"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"text/template"
	"time"
)

type Templates struct {
	fs   fauxfile.Filesystem
	log  *log.Logger
	root *template.Template
	fmap *template.FuncMap
}

func NewTemplates() *Templates {
	return &Templates{
		fs:   &fauxfile.RealFilesystem{},
		log:  log.New(os.Stderr, "", log.LstdFlags),
		root: template.New("root"),
		fmap: getFuncMap(),
	}
}

// Includes the contents of the supplied directory in the root template.
func (ts *Templates) AddTemplatesFromDir(path string) (err error) {
	var (
		names []string
		fpath string
	)
	if names, err = ts.readDir(path); err != nil {
		ts.log.Printf("Templates directory not found %v\n", path)
		return
	}
	for _, n := range names {
		fpath = filepath.Join(path, n)
		ts.log.Printf("Adding template at %v\n", fpath)
		if err = ts.AddTemplateFromFile(fpath); err != nil {
			return
		}
	}
	return
}

// Includes the contents of the supplied file in the root template.
func (ts *Templates) AddTemplateFromFile(path string) (err error) {
	var (
		text string
	)
	if text, err = ts.readFile(path); err != nil {
		return
	}
	err = ts.AddTemplate(text)
	return
}

// Includes the contents of the supplied text in the root template.
func (ts *Templates) AddTemplate(text string) (err error) {
	_, err = ts.root.Funcs(*ts.fmap).Parse(text)
	return
}

// Overrides portions of the root template with a file's contents and renders.
func (ts *Templates) RenderFile(path string, data map[string]interface{}) (out string, err error) {
	var text string
	if text, err = ts.readFile(path); err != nil {
		return
	}
	out, err = ts.RenderText(text, data)
	return
}

// Overrides portions of the root template and renders the appropriate data.
func (ts *Templates) RenderText(text string, data map[string]interface{}) (out string, err error) {
	var (
		clone  *template.Template
		tmpl   *template.Template
		writer *bytes.Buffer
	)
	if tmpl, err = template.New("temp").Funcs(*ts.fmap).Parse(text); err != nil {
		return
	}
	if clone, err = ts.mergeTemplate(tmpl); err != nil {
		return
	}
	writer = bytes.NewBufferString("")
	if err = clone.Execute(writer, data); err == nil {
		out = writer.String()
	}
	return
}

// Renders the root template without any overrides.
func (ts *Templates) Render(data map[string]interface{}) (out string, err error) {
	writer := bytes.NewBufferString("")
	if err = ts.root.Execute(writer, data); err == nil {
		out = writer.String()
	}
	return
}

// Reads directory contents from the given path and returns file names.
func (ts *Templates) readDir(path string) (names []string, err error) {
	var f fauxfile.File
	if f, err = ts.fs.Open(path); err != nil {
		return
	}
	defer f.Close()
	names, err = f.Readdirnames(-1)
	return
}

// Returns a copy of the root template with the supplied template merged in.
func (ts *Templates) mergeTemplate(t *template.Template) (out *template.Template, err error) {
	defer func() {
		if r := recover(); r != nil {
			// Seems to be a bug with cloning empty templates.
			err = fmt.Errorf("Problem cloning template: %v", r)
		}
	}()
	if out, err = ts.root.Clone(); err != nil {
		return
	}
	for _, tmpl := range t.Templates() {
		ptr := out.Lookup(tmpl.Name())
		if ptr == nil {
			ptr = out.New(tmpl.Name())
		}
		(*ptr) = *tmpl
	}
	return
}

// Reads a file from the given path and returns a string of the contents.
func (ts *Templates) readFile(path string) (out string, err error) {
	var (
		f   fauxfile.File
		fi  os.FileInfo
		buf []byte
	)
	if f, err = ts.fs.Open(path); err != nil {
		return
	}
	defer f.Close()
	if fi, err = f.Stat(); err != nil {
		return
	}
	buf = make([]byte, fi.Size())
	if _, err = f.Read(buf); err != nil {
		if err != io.EOF {
			return
		}
		err = nil
	}
	out = string(buf)
	return
}

// Returns a base set of functions for use in templates.
func getFuncMap() *template.FuncMap {
	return &template.FuncMap{
		"timeformat": func(t time.Time, f string) string {
			var format_str string
			switch f {
			case "ANSIC":
				format_str = time.ANSIC
			case "UnixDate":
				format_str = time.UnixDate
			case "RubyDate":
				format_str = time.RubyDate
			case "RFC822":
				format_str = time.RFC822
			case "RFC822Z":
				format_str = time.RFC822Z
			case "RFC850":
				format_str = time.RFC850
			case "RFC1123":
				format_str = time.RFC1123
			case "RFC1123Z":
				format_str = time.RFC1123Z
			case "RFC3339":
				format_str = time.RFC3339
			case "RFC3339Nano":
				format_str = time.RFC3339Nano
			case "Kitchen":
				format_str = time.Kitchen
			case "Stamp":
				format_str = time.Stamp
			case "StampMilli":
				format_str = time.StampMilli
			case "StampMicro":
				format_str = time.StampMicro
			case "StampNano":
				format_str = time.StampNano
			default:
				format_str = f
			}
			return t.Format(format_str)
		},
		"textcontent": func(s string) string {
			rex, _ := regexp.Compile("<[^>]*>")
			return rex.ReplaceAllLiteralString(s, "")
		},
	}
}
