package main

import (
	"html/template"
	"io"
	"path"

	"github.com/Masterminds/sprig"
)

// RenderTemplate renders Team using tplFile into target Writer
func RenderTemplate(t *Team, tplFile string, wr io.Writer) error {
	tpl, err := template.New(path.Base(tplFile)).
		Funcs(sprig.FuncMap()).
		ParseFiles(tplFile)
	if err != nil {
		return err
	}
	return tpl.Execute(wr, t)
}
