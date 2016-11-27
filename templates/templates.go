package templates

import (
	"fmt"
	"github.com/boreq/blogs/http/context"
	"github.com/boreq/blogs/logging"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var log = logging.GetLogger("templates")
var templates map[string]*template.Template

func Load() error {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	templatesDir := "_templates/"

	var layouts []string
	layoutsDir := templatesDir + "templates/"
	filepath.Walk(layoutsDir,
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				layouts = append(layouts, path)
			}
			return nil
		})

	var snippets []string
	snippetsDir := templatesDir + "snippets/"
	filepath.Walk(snippetsDir,
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				snippets = append(snippets, path)
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
		files = append(files, snippets...)
		for _, fname := range files {
			log.Debugf("    dependency %s", fname)
		}
		key := layout[len(layoutsDir):]
		templates[key] = template.Must(template.ParseFiles(files...))
	}
	return nil
}

func GetDefaultData(r *http.Request) map[string]interface{} {
	var data = make(map[string]interface{})
	var ctx = context.Get(r)
	data["user"] = ctx.User
	data["now"] = time.Now()
	return data
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
