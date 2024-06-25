package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	required = "REQUIRED"
	optional = "OPTIONAL"
	web      = "web"
	pdf      = "pdf"
)

func main() {
	var inpath = flag.String("in", required, "path to an input markdown file")
	var outpath = flag.String("out", optional, "path to an output html file")
	var title = flag.String("title", optional, "title for the generated html doc")
	var formatFor = flag.String("format-for", web, fmt.Sprintf("will the rendered markdown be used for %s or %s?", web, pdf))
	flag.Parse()

	if *inpath == "" || *inpath == required {
		checkErr(errors.New("in is required"))
	}
	if *outpath == "" || *outpath == optional {
		*outpath = fileName(*inpath) + ".html"
	}
	if *title == "" || *title == optional {
		*title = fileName(*inpath)
	}
	if *formatFor != web && *formatFor != pdf {
		checkErr(fmt.Errorf("type must be %s or %s", web, pdf))
	}

	markdown, err := ioutil.ReadFile(*inpath)
	checkErr(err)

	rendered, err := renderMarkdown(markdown)
	checkErr(err)

	var builder docBuilder
	switch *formatFor {
	case web:
		builder = newWebDoc(*title, string(rendered))
	case pdf:
		builder = newPDFDoc(*title, string(rendered))
	}

	generated, err := builder.buildDoc()
	checkErr(err)

	if *outpath == optional {
		fmt.Println(generated)
	} else {
		outfile, err := os.Create(*outpath)
		checkErr(err)
		_, err = outfile.WriteString(generated)
		checkErr(err)
	}
}

func renderMarkdown(markdown []byte) ([]byte, error) {
	body := map[string]string{
		"text": string(markdown),
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := http.Post(
		"https://api.github.com/markdown",
		"application/json",
		bytes.NewReader(bodyJSON),
	)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func fileName(path string) string {
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type docBuilder struct {
	Title      string
	FontSize   string
	LineHeight string
	Body       string
}

func (builder docBuilder) buildDoc() (string, error) {
	docTempl, err := template.New("doc").Parse(docTemplate)
	if err != nil {
		return "", err
	}
	var docBuf bytes.Buffer
	err = docTempl.Execute(&docBuf, builder)
	if err != nil {
		return "", err
	}
	return docBuf.String(), nil
}

func newWebDoc(title string, body string) docBuilder {
	return docBuilder{
		title,
		"1.1em",
		"1.5em",
		body,
	}
}

func newPDFDoc(title string, body string) docBuilder {
	return docBuilder{
		title,
		"0.9em",
		"1.3em",
		body,
	}
}

var docTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
    <meta name="viewport" content="width=device-width">
    <title>{{.Title}}</title>
    <style type="text/css">
        body {
            font-size: {{.FontSize}};
            line-height: {{.LineHeight}};
            max-width: 45em;
            margin: auto;
            padding: 0 2%;
            font-family: system-ui;
        }
        img {
            max-width: 100%;
            margin-left: 0.5em;
            vertical-align: -0.3em;
        }
    </style>
</head>
<body>
{{.Body}}
</body>
`
