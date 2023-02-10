package wecom

import (
	"embed"
	"github.com/AkvicorEdwards/glog"
	"html/template"
)

//go:embed static/html/*
var htmlFS embed.FS

var tpls *template.Template

var tplIndex *template.Template
var tplSecret *template.Template
var tplSecretInsert *template.Template

func init() {
	tpls = template.Must(template.ParseFS(htmlFS, "static/html/*"))

	tplIndex = tpls.Lookup("index.gohtml")
	if tplIndex == nil {
		glog.Fatal("missing html template [index.gohtml]")
	}
	tplSecret = tpls.Lookup("secret.gohtml")
	if tplSecret == nil {
		glog.Fatal("missing html template [secret.gohtml]")
	}
	tplSecretInsert = tpls.Lookup("secret_insert.gohtml")
	if tplSecretInsert == nil {
		glog.Fatal("missing html template [secret_insert.gohtml]")
	}
}
