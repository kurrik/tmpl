tmpl
====
This library offers a better approach toward the management of Go's built-in
template libraries.

The goal is to be able to configure a set of root templates, presumably for
a website, and then be able to render overrides as needed.

Installing
----------
Run

    go get github.com/kurrik/tmpl

Include in your source:

    import "github.com/kurrik/tmpl"

Godoc
-----
See http://godoc.org/github.com/kurrik/tmpl

Testing
-------
In the project root run:

    go test

Using
-----

===Layout===

Start by defining a root template which contains placeholders for all
referenced sub-templates.  In the example, one file is used to hold all these
but they could easily be split out into multiple files.

**templates/root/root.tmpl**

    <html>
      <head>{{template "head" .}}</head>
      <body>{{template "body" .}}</body>
    </html>
    {{define "head"}}<title>{{.Title}}</title>{{end}}
    {{define "body"}}{{end}}

Then define overrides which will be selectively used.  For example, here is
a post template which overrides body:

**templates/post.tmpl**

    {{define "body"}}
    <div>
      <h2>{{.Title}}</h2>
      <p>{{timeformat .Date "UnixDate"}}</p>
      <p>{{.Body}}</p>
    </div>
    {{end}}

Note that it does not reside in the same directory as the root template, this
allows us to parse the root template directory in a single command, and
selectively use the post template override as needed.

Here's another body override which renders a list of posts:

**templates/index.tmpl**

    {{define "body"}}
    {{range .Posts}}
    <div>
      <h2>{{.Title}}</h2>
      <p>{{timeformat .Date "UnixDate"}}</p>
      <p>{{textcontent .Body}}</p>
    </div>
    {{end}}
    {{end}}

===Rendering===

To render templates call `tmpl.NewTemplates` and then parse the root template
structure.  Since the example organized the root template into a directory,
this is as simple as calling `AddTemplatesFromDir`.

	templates := tmpl.NewTemplates()
	if err = templates.AddTemplatesFromDir("templates/root"); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

Set up some sample data to render:

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

Rendering templates is done by calling any of the `RenderXXXX` methods,
passing in override templates as needed.  Here is the call to render the
index list of posts:

	// Render the index.
	if out, err = templates.RenderFile("templates/index.tmpl", posts); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Index\n=============================")
		fmt.Println(out)
	}

This produces:

    Index
    =============================
    <html>
      <head><title>My Site</title></head>
      <body>

    <div>
      <h2>Hello World</h2>
      <p>Tue Nov 10 23:00:00 UTC 2009</p>
      <p>Hi!</p>
    </div>

    <div>
      <h2>Hello Again</h2>
      <p>Wed Nov 11 00:00:00 UTC 2009</p>
      <p>Hi there!</p>
    </div>

    </body>
    </html>

Here is the call to render the first post:

	// Render a post.
	if out, err = templates.RenderFile("templates/post.tmpl", post1); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Post\n=============================")
		fmt.Println(out)
	}

This produces:

    Post
    =============================
    <html>
      <head><title>Hello World</title></head>
      <body>
    <div>
      <h2>Hello World</h2>
      <p>Tue Nov 10 23:00:00 UTC 2009</p>
      <p><em>Hi!</em></p>
    </div>
    </body>
    </html>

You may also render a template string with the `RenderText` method:

	// Render an ad-hoc template.
	if out, err = templates.RenderText(`{{define "body"}}{{.Title}}{{end}}`, post1); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Ad-hoc\n=============================")
		fmt.Println(out)
	}

This produces:

    Ad-hoc
    =============================
    <html>
      <head><title>Hello World</title></head>
      <body>Hello World</body>
    </html>

===Functions===

Some template functions are provided:

**timeformat**
This function takes a `time.Time` value as its first argument and renders
it according to a format string supplied as its second argument.  If a string
is passed which matches a predefined format defined in the `time` package, that
format is used, otherwise it is treated as a literal format string:

    {{timeformat .Date "UnixDate"}}
    {{timeformat .Date "Mon Jan 02 15:04:05 -0700 2006"}}

**textcontent**
This performs rudmentary stripping of HTML tags so that only text content
is output.  IMPORTANT! Do not use this method as a security measure - proper
sanitization should be performed on untrusted input.

    {{textcontent .Body}}
