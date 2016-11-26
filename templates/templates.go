package templates

import (
	"fmt"
	"github.com/boreq/blogs/logging"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var log = logging.GetLogger("templates")
var templates map[string]*template.Template

func Load() error {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	templatesDir := "_templates/"

	var layouts []string
	filepath.Walk(templatesDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			layouts = append(layouts, path)
		}
		return nil
	})

	for _, layout := range layouts {
		log.Debugf("Loading %s", layout)
		parentDir := filepath.Clean(filepath.Dir(layout) + "/..")
		if parentDir == "." {
			parentDir = ""
		}

		files := make([]string, 0)
		for _, fileName := range layouts {
			fileDir := filepath.Clean(filepath.Dir(fileName))
			if strings.HasPrefix(parentDir, fileDir) {
				files = append(files, fileName)
			}
		}
		files = append(files, layout)
		for _, fname := range files {
			log.Debugf("    dependency %s", fname)
		}
		key := layout[len(templatesDir):]
		templates[key] = template.Must(template.ParseFiles(files...))
	}
	return nil
}

// renderTemplate is a wrapper around template.ExecuteTemplate.
func RenderTemplate(w http.ResponseWriter, name string, data map[string]interface{}) error {
	tmpl, ok := templates[name]
	if !ok {
		return fmt.Errorf("The template %s does not exist.", name)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.ExecuteTemplate(w, "base", data)
}
