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
	"strings"
	"testing"
	"text/template"
	"time"
)

func LooseCompare(t *testing.T, a string, b string) bool {
	a = strings.Replace(a, " ", "", -1)
	a = strings.Replace(a, "\t", "", -1)
	a = strings.Replace(a, "\n", "", -1)
	b = strings.Replace(b, " ", "", -1)
	b = strings.Replace(b, "\t", "", -1)
	b = strings.Replace(b, "\n", "", -1)
	if a != b {
		t.Logf("LooseCompare diff:\n%v\n%v", a, b)
		return false
	}
	return true
}

const (
	TMPL_BASE = `
		[h]{{template "head" .}}[/h]
		[b]{{template "body" .}}[/b]
	`
	TMPL_HEAD = `
		{{define "head"}}
			[c]{{.HeadContent}}[/c]
		{{end}}
	`
	TMPL_BODY = `
		{{define "body"}}
			[c]{{.BodyContent}}[/c]
		{{end}}
	`
	TMPL_WRAP = `
		{{define "wrap"}}
			[w]{{template "wrap_content" .}}[/w]
		{{end}}
		{{define "wrap_content"}}{{end}}
	`
	TMPL_CONT = `
		{{define "body"}}
			{{template "wrap" .}}
		{{end}}
		{{define "wrap_content"}}[c]{{.BodyContent}}[/c]{{end}}
	`
)

func TestRender(t *testing.T) {
	var (
		out  string
		err  error
		data = map[string]interface{}{
			"HeadContent": "hc",
			"BodyContent": "bc",
		}
	)

	templates := NewTemplates()
	templates.AddTemplate(TMPL_BASE)
	templates.AddTemplate(TMPL_HEAD)
	templates.AddTemplate(TMPL_BODY)

	if out, err = templates.Render(data); err != nil {
		t.Fatalf("Error rendering: %v", err)
	}
	if !LooseCompare(t, out, "[h][c]hc[/c][/h][b][c]bc[/c][/b]") {
		t.Fatalf("Simple render did not produce correct output")
	}
}

func TestWrappedContent(t *testing.T) {
	var (
		out  string
		err  error
		data = map[string]interface{}{
			"HeadContent": "hc",
			"BodyContent": "bc",
		}
	)

	templates := NewTemplates()
	templates.AddTemplate(TMPL_BASE)
	templates.AddTemplate(TMPL_HEAD)
	templates.AddTemplate(TMPL_BODY)
	templates.AddTemplate(TMPL_WRAP)

	if out, err = templates.RenderText(TMPL_CONT, data); err != nil {
		t.Fatalf("Error rendering: %v", err)
	}
	if !LooseCompare(t, out, "[h][c]hc[/c][/h][b][w][c]bc[/c][/w][/b]") {
		t.Fatalf("TestWrappedContent error")
	}
}

func TestRenderText(t *testing.T) {
	var (
		out  string
		err  error
		tmpl = `{{define "body"}}[BC]{{.BodyContent}}[/BC]{{end}}`
		data = map[string]interface{}{
			"HeadContent": "hc",
			"BodyContent": "bc",
		}
	)

	templates := NewTemplates()
	templates.AddTemplate(TMPL_BASE)
	templates.AddTemplate(TMPL_HEAD)
	templates.AddTemplate(TMPL_BODY)

	if out, err = templates.RenderText(tmpl, data); err != nil {
		t.Fatalf("Error rendering: %v", err)
	}
	if !LooseCompare(t, out, "[h][c]hc[/c][/h][b][BC]bc[/BC][/b]") {
		t.Fatalf("TestRenderText error")
	}
}

func TestRenderTemplate(t *testing.T) {
	var (
		out  string
		err  error
		tmpl *template.Template
		data = map[string]interface{}{
			"HeadContent": "hc",
			"BodyContent": "bc",
		}
	)
	tmpl, _ = template.New("body").Parse(`[BC]{{.BodyContent}}[/BC]`)
	templates := NewTemplates()
	templates.AddTemplate(TMPL_BASE)
	templates.AddTemplate(TMPL_HEAD)
	templates.AddTemplate(TMPL_BODY)

	if out, err = templates.RenderTemplate(tmpl, data); err != nil {
		t.Fatalf("Error rendering: %v", err)
	}
	if !LooseCompare(t, out, "[h][c]hc[/c][/h][b][BC]bc[/BC][/b]") {
		t.Fatalf("TestRenderTemplate error")
	}
}

func TestOutOfOrderTemplateInitialization(t *testing.T) {
	var (
		out  string
		err  error
		data = map[string]interface{}{
			"HeadContent": "hc",
			"BodyContent": "bc",
		}
	)

	templates := NewTemplates()
	templates.AddTemplate(TMPL_BODY)
	templates.AddTemplate(TMPL_HEAD)
	templates.AddTemplate(TMPL_BASE)

	if out, err = templates.Render(data); err != nil {
		t.Fatalf("Error rendering: %v", err)
	}
	if !LooseCompare(t, out, "[h][c]hc[/c][/h][b][c]bc[/c][/b]") {
		t.Fatalf("TestOutOfOrderTemplateInitialization error")
	}
}

func TestAddTemplateFromTemplate(t *testing.T) {
	var (
		out  string
		err  error
		tmpl *template.Template
		data = map[string]interface{}{
			"HeadContent": "hc",
			"BodyContent": "bc",
		}
	)

	templates := NewTemplates()
	templates.AddTemplate(TMPL_BASE)
	templates.AddTemplate(TMPL_HEAD)
	tmpl, _ = template.New("body").Parse(`[BC]{{.BodyContent}}[/BC]`)
	if err = templates.AddTemplateFromTemplate(tmpl); err != nil {
		t.Fatalf("Error adding template: %v", err)
	}

	if out, err = templates.Render(data); err != nil {
		t.Fatalf("Error rendering: %v", err)
	}
	if !LooseCompare(t, out, "[h][c]hc[/c][/h][b][BC]bc[/BC][/b]") {
		t.Fatalf("TestAddTemplateFromTemplate error")
	}
}

func TestOutOfOrderAddTemplateFromTemplate(t *testing.T) {
	var (
		out  string
		err  error
		tmpl *template.Template
		data = map[string]interface{}{
			"HeadContent": "hc",
			"BodyContent": "bc",
		}
	)

	templates := NewTemplates()
	tmpl, _ = template.New("body").Parse(`[BC]{{.BodyContent}}[/BC]`)
	if err = templates.AddTemplateFromTemplate(tmpl); err != nil {
		t.Fatalf("Error adding template: %v", err)
	}
	templates.AddTemplate(TMPL_BASE)
	templates.AddTemplate(TMPL_HEAD)

	if out, err = templates.Render(data); err != nil {
		t.Fatalf("Error rendering: %v", err)
	}
	if !LooseCompare(t, out, "[h][c]hc[/c][/h][b][BC]bc[/BC][/b]") {
		t.Fatalf("TestOutOfOrderAddTemplateFromTemplate error")
	}
}

func TestCallTemplateFromUpdatedTemplate(t *testing.T) {
	var (
		out  string
		err  error
		data = map[string]interface{}{}
	)
	templates := NewTemplates()
	templates.AddTemplate(`{{define "root"}}[r]{{template "body" .}}[r]{{end}}`)
	templates.AddTemplate(`{{define "body"}}{{end}}`)
	templates.AddTemplate(`{{define "helper"}}Helper{{end}}`)
	// Update body and call a loaded template from it.
	if err = templates.AddTemplate(`{{define "body"}}{{template "helper" .}}{{end}}`); err != nil {
		t.Fatalf("Could not update template: %v", err)
	}
	if out, err = templates.Render(data); err != nil {
		t.Fatalf("Error rendering: %v", err)
	}
	if !LooseCompare(t, out, "[r]Helper[r]") {
		t.Fatalf("TestCallTemplateFromUpdatedTemplate error")
	}
}

func TestNamedRenderTemplate(t *testing.T) {
	var (
		out  string
		err  error
		data = map[string]interface{}{}
	)
	templates := NewTemplates()
	templates.AddTemplate(`{{define "root"}}[r]{{template "body" .}}[r]{{end}}`)
	templates.AddTemplate(`{{define "body"}}{{end}}`)
	templates.AddTemplate(`{{define "helper"}}Helper{{end}}`)
	// Update body and call a loaded template from it.
	if err = templates.AddTemplate(`{{define "body"}}{{template "helper" .}}{{end}}`); err != nil {
		t.Fatalf("Could not update template: %v", err)
	}
	if out, err = templates.NamedRender("body", data); err != nil {
		t.Fatalf("Error rendering: %v", err)
	}
	if !LooseCompare(t, out, "Helper") {
		t.Fatalf("TestNamedRenderTemplate error")
	}
}

func TestTextcontentFunction(t *testing.T) {
	var (
		out  string
		err  error
		tmpl = `{{define "body"}}[TC]{{textcontent .BodyContent}}[/TC]{{end}}`
		data = map[string]interface{}{
			"HeadContent": "hc",
			"BodyContent": "<a>bc</a>",
		}
	)

	templates := NewTemplates()
	templates.AddTemplate(TMPL_BODY)
	templates.AddTemplate(TMPL_HEAD)
	templates.AddTemplate(TMPL_BASE)

	if out, err = templates.RenderText(tmpl, data); err != nil {
		t.Fatalf("Error rendering: %v", err)
	}
	if !LooseCompare(t, out, "[h][c]hc[/c][/h][b][TC]bc[/TC][/b]") {
		t.Fatalf("TestTextcontentFunction error")
	}
}

func TestTimeformatFunction(t *testing.T) {
	var (
		out  string
		err  error
		tmpl = `{{define "body"}}[TF]{{timeformat .BodyContent "UnixDate"}}[/TF]{{end}}`
		data = map[string]interface{}{
			"HeadContent": "hc",
			"BodyContent": time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		}
	)

	templates := NewTemplates()
	templates.AddTemplate(TMPL_BODY)
	templates.AddTemplate(TMPL_HEAD)
	templates.AddTemplate(TMPL_BASE)

	if out, err = templates.RenderText(tmpl, data); err != nil {
		t.Fatalf("Error rendering: %v", err)
	}
	if !LooseCompare(t, out, "[h][c]hc[/c][/h][b][TF]TueNov1023:00:00UTC2009[/TF][/b]") {
		t.Fatalf("TestTimeformatFunction error")
	}
}
