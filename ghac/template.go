package main

import (
	"html/template"
	"io"
	"path"
)

// Action interface for each kind of struct
type Action interface {
	RenderTemplate(tplFile string, wr io.Writer) error
}

// RenderTemplate renders Team using tplFile into target Writer
func (t Team) RenderTemplate(tplFile string, wr io.Writer) error {
	tpl, err := template.New(path.Base(tplFile)).
		ParseFiles(tplFile)
	if err != nil {
		return err
	}
	return tpl.Execute(wr, t)
}

// RenderTemplate renders TeamVault using tplFile into target Writer
func (tv TeamVault) RenderTemplate(tplFile string, wr io.Writer) error {
	tpl, err := template.New(path.Base(tplFile)).
		ParseFiles(tplFile)
	if err != nil {
		return err
	}
	return tpl.Execute(wr, tv)
}
