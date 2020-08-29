package generator

import (
	"bytes"
	"go/format"
	"io"
	"text/template"
)

const (
	templatePath = "go-validate/generator/templates/package.tmpl"
)

func Render(w io.Writer, data *TemplateData) error {
	t, err := template.ParseFiles(
		"go-validate/generator/templates/package.tmpl",
		"go-validate/generator/templates/validators.tmpl",
	)
	if err != nil {
		return err
	}

	// compile template
	template.Must(t, err)

	// render template to buffer
	var renderedBuf bytes.Buffer

	err = t.Execute(&renderedBuf, data)
	if err != nil {
		return err
	}

	// format as gofmt
	formattedBuf, err := format.Source(renderedBuf.Bytes())
	if err != nil {
		return err
	}

	_, err = w.Write(formattedBuf)

	return err
}
