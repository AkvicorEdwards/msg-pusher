package tpl

import (
	"embed"
	_ "embed"
	"github.com/AkvicorEdwards/glog"
	"html/template"
)

//go:embed static/img/favicon.ico
var Favicon []byte

//go:embed static/html/*
var html embed.FS

var Login *template.Template
var Index *template.Template
var Secret *template.Template
var SecretInsert *template.Template
var Target *template.Template

func init() {
	t := template.Must(template.ParseFS(html, "static/html/*"))

	Login = t.Lookup("login.gohtml")
	if Login == nil {
		glog.Fatal("missing html template [login.gohtml]")
	}
	Index = t.Lookup("index.gohtml")
	if Index == nil {
		glog.Fatal("missing html template [index.gohtml]")
	}
	Secret = t.Lookup("secret.gohtml")
	if Secret == nil {
		glog.Fatal("missing html template [secret.gohtml]")
	}
	SecretInsert = t.Lookup("secret_insert.gohtml")
	if SecretInsert == nil {
		glog.Fatal("missing html template [secret_insert.gohtml]")
	}
	Target = t.Lookup("target.gohtml")
	if Target == nil {
		glog.Fatal("missing html template [target.gohtml]")
	}
}
