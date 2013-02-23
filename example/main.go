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

package main

import (
	"../" // Use github.com/kurrik/tmpl for your code.
	"fmt"
	"time"
)

func main() {
	var (
		out   string
		err   error
	)

	post1 := map[string]interface{}{
		"Title": "Hello World",
		"Date":  time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
		"Body":  "<em>Hi!</em>",
	}

	post2 := map[string]interface{}{
		"Title": "Hello Again",
		"Date":  time.Date(2009, time.November, 10, 24, 0, 0, 0, time.UTC),
		"Body":  "<em>Hi there!</em>",
	}

	posts := map[string]interface{}{
		"Title": "My Site",
		"Posts": []map[string]interface{}{
			post1,
			post2,
		},
	}

	templates := tmpl.NewTemplates()
	if err = templates.AddTemplatesFromDir("templates/root"); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Render the root template without any overrides.
	if out, err = templates.Render(posts); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Root\n=============================")
		fmt.Println(out)
	}

	// Render the index.
	if out, err = templates.RenderFile("templates/index.tmpl", posts); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Index\n=============================")
		fmt.Println(out)
	}

	// Render a post.
	if out, err = templates.RenderFile("templates/post.tmpl", post1); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Post\n=============================")
		fmt.Println(out)
	}

	// Render an ad-hoc template.
	if out, err = templates.RenderText(`{{define "body"}}{{.Title}}{{end}}`, post1); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Ad-hoc\n=============================")
		fmt.Println(out)
	}
}
