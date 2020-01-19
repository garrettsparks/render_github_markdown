# render_github_markdown
tool for rendering markdown as html for hosting as a webpage or printing to pdf

utilizes the [github markdown API](https://developer.github.com/v3/markdown/) for markdown to html conversion

### usage
```
  -format-for string
        will the rendered markdown be used for web or pdf? (default "web")
  -in string
        path to an input markdown file (default "REQUIRED")
  -out string
        path to an output html file (default is the name of the source file with the extension .html)
  -title string
        title for the generated html doc (default is the name of "in" without the extension)
```

### examples
```
$ ./render_github_markdown --in README.md
```
outputs `README.html` like [this](https://garrettsparks.github.io/README)

```
$ ./render_github_markdown --in README.md --out PDF_README.html --title "PDF ReadMe" --format-for pdf
```
outputs `PDF_README.html` like [this](https://garrettsparks.github.io/PDF_README) which is slightly better formatted for converting to PDF
