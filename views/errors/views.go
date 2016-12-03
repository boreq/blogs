package errors

import (
	"bytes"
	"github.com/boreq/blogs/templates"
	"net/http"
)

func displayErrorPage(code int, message string, w http.ResponseWriter, r *http.Request) {
	var data = templates.GetDefaultData(r)
	data["error_code"] = code
	data["error_message"] = message
	err := templates.RenderTemplateSafe(w, "errors/error.tmpl", data)
	if err != nil {
		InternalServerError(w, r)
	}
}

func BadRequest(w http.ResponseWriter, r *http.Request) {
	displayErrorPage(400, "Bad Request", w, r)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	displayErrorPage(404, "Not Found", w, r)
}

var internalServerErrorResponse = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<title>Internal server error</title>
		<style>
			body {
				background: #fff;
				color: #000;
				font-family: sans-serif;
				margin: 200px auto 0 auto;
				width: 500px;
			}
			a {
				text-decoration: underline;
				color: inherit;
			}
			a:hover {
				text-decoration: none;
			}
			.h1{font-size: 25px;}
			.p{font-size: 15px;}
		</style>
<body>
	<h1>Internal server error</h1>
	<p>Something went really wrong this time.</p>
	<p><a href="/">Homepage</a></p>
`

func InternalServerError(w http.ResponseWriter, r *http.Request) {
	buf := bytes.NewBufferString(internalServerErrorResponse)
	w.WriteHeader(500)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}
