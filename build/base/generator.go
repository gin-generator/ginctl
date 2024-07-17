package base

import (
	"sync"
	"text/template"
)

var (
	TemplateCache sync.Map
)

func GetTemplateParse(filePath string) (*template.Template, error) {
	if t, ok := TemplateCache.Load(filePath); ok {
		return t.(*template.Template), nil
	}

	// The template is parsed and cached
	t, err := template.ParseFiles(filePath)
	if err != nil {
		return nil, err
	}

	TemplateCache.Store(filePath, t)
	return t, nil
}
