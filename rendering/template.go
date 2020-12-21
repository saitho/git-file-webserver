package rendering

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"strings"

	"github.com/markbates/pkger"
)

func RenderTemplate(path string, params interface{}) (string, error) {
	tplFuncMap := make(template.FuncMap)
	tplFuncMap["Split"] = strings.Split

	f, _ := pkger.Open(path)
	defer f.Close()
	sl, _ := ioutil.ReadAll(f)

	t, err := template.New("").Funcs(tplFuncMap).Parse(string(sl))
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, params); err != nil {
		return "", err
	}
	return tpl.String(), nil
}
