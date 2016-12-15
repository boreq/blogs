// Package templates loads and renders HTML templates.
package templates

import (
	"bytes"
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

// Load should be called before rendering templates using functions contained
// in this package.
func Load(templatesDir string) error {
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		return err
	}

	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	templatesDir = strings.TrimRight(templatesDir, string(os.PathSeparator))
	layoutsDir := templatesDir + "/templates/"
	snippetsDir := templatesDir + "/snippets/"
	var layouts = findFiles(layoutsDir)
	var snippets = findFiles(snippetsDir)

	for _, layout := range layouts {
		log.Debugf("Loading %s", layout)
		files := make([]string, 0)

		// Files from parent directories
		parentDir := filepath.Clean(filepath.Dir(layout) + "/..")
		if parentDir == "." {
			parentDir = ""
		}
		for _, fileName := range layouts {
			fileDir := filepath.Clean(filepath.Dir(fileName))
			if strings.HasPrefix(parentDir, fileDir) {
				files = append(files, fileName)
			}
		}

		// Partial templates
		for _, fileName := range layouts {
			basePath := fileName[:len(fileName)-len(filepath.Ext(fileName))]
			if strings.HasPrefix(layout, basePath) && layout[len(basePath)] == '_' {
				files = append(files, fileName)
			}
		}

		for _, fname := range files {
			log.Debugf("    dependency %s", fname)
		}
		log.Debugf("    +%d snippets", len(snippets))

		// Layout and snippets
		files = append(files, layout)
		files = append(files, snippets...)

		key := layout[len(layoutsDir):]
		templates[key] = template.Must(template.New("").Funcs(getFuncMap()).ParseFiles(files...))
	}
	return nil
}

func findFiles(directory string) []string {
	var rv []string
	filepath.Walk(directory,
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				rv = append(rv, path)
			}
			return nil
		})
	return rv
}

// GetDefaultData returns the base data which can be extended with new keys
// and passed to the RenderTemplate or RenderTemplateSafe functions. It is
// advised to use this function to get the initial map instead of creating it
// directly.
func GetDefaultData(r *http.Request) map[string]interface{} {
	var data = make(map[string]interface{})
	var ctx = context.Get(r)
	data["user"] = ctx.User
	data["request"] = r
	data["now"] = time.Now()
	return data
}

func getTemplate(name string) (*template.Template, error) {
	tmpl, ok := templates[name]
	if !ok {
		return nil, fmt.Errorf("The template %s does not exist.", name)
	}
	return tmpl, nil
}

func setContentTypeHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

// RenderTemplate is a wrapper around template.ExecuteTemplate. If an error
// occurs in this function partial results may be written to the output writer.
func RenderTemplate(w http.ResponseWriter, name string, data map[string]interface{}) error {
	tmpl, err := getTemplate(name)
	if err != nil {
		return err
	}

	setContentTypeHeader(w)
	return tmpl.ExecuteTemplate(w, "base", data)
}

// RenderTemplateSafe is a wrapper around template.ExecuteTemplate. If an error
// occurs this function will attempt to ensure that no data will written to
// the output writer.
func RenderTemplateSafe(w http.ResponseWriter, name string, data map[string]interface{}) error {
	tmpl, err := getTemplate(name)
	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	err = tmpl.ExecuteTemplate(buf, "base", data)
	if err != nil {
		return err
	}

	setContentTypeHeader(w)
	_, err = buf.WriteTo(w)
	return err
}
