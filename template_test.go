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
	"time"
)

func LooseCompare(t *testing.T, a string, b string) bool {
	a = strings.Replace(a, " ", "", -1)
	a = strings.Replace(a, "\n", "", -1)
	b = strings.Replace(b, " ", "", -1)
	b = strings.Replace(b, "\n", "", -1)
	if a != b {
		t.Logf("LooseCompare diff:\n%v\n%v", a, b)
		return false
	}
	return true
}

const (
	TMPL_BASE = `[h]{{template "head" .}}[/h][b]{{template "body" .}}[/b]`
	TMPL_HEAD = `{{define "head"}}[c]{{.HeadContent}}[/c]{{end}}`
	TMPL_BODY = `{{define "body"}}[c]{{.BodyContent}}[/c]{{end}}`
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
		t.Fatalf("RenderText did not produce correct output")
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
		t.Fatalf("Out of order init did not produce correct output")
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
		t.Fatalf("Textcontent function did not produce correct output")
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
		t.Fatalf("Timeformat function did not produce correct output")
	}
}



