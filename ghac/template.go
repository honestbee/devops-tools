package main

import (
	"html/template"
	"io"
	"path"
)

// RenderTemplate renders Team using tplFile into target Writer
func RenderTemplate(t Team, tplFile string, wr io.Writer) error {
	tpl, err := template.New(path.Base(tplFile)).
		ParseFiles(tplFile)
	if err != nil {
		return err
	}
	return tpl.Execute(wr, t)
}

// RenderCustomTemplate renders TeamVault using tplFile into target Writer
func RenderCustomTemplate(m map[string]string, tplFile string, wr io.Writer) error {
	tpl, err := template.New(path.Base(tplFile)).
		ParseFiles(tplFile)
	if err != nil {
		return err
	}
	return tpl.Execute(wr, m)
}
