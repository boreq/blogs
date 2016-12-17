package errors

import (
	"bytes"
	"fmt"
	"github.com/boreq/blogs/templates"
	"net/http"
	"runtime/debug"
)

func BadRequest(w http.ResponseWriter, r *http.Request) {
	displayErrorPageOrInternalServerError(400, "Bad Request", w, r)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	displayErrorPageOrInternalServerError(404, "Not Found", w, r)
}

func InternalServerErrorWithStack(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Printf("%s\n", err)
	fmt.Println(string(debug.Stack()))
	InternalServerError(w, r)
}

func InternalServerError(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			buf := bytes.NewBufferString(internalServerErrorResponse)
			w.WriteHeader(500)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			buf.WriteTo(w)
		}
		return

	}()
	err := displayErrorPage(500, "Internal Server Error", w, r)
	if err != nil {
		panic(err)
	}
}

func displayErrorPageOrInternalServerError(code int, message string, w http.ResponseWriter, r *http.Request) {
	if err := displayErrorPage(code, message, w, r); err != nil {
		InternalServerError(w, r)
	}
}

func displayErrorPage(code int, message string, w http.ResponseWriter, r *http.Request) error {
	var data = templates.GetDefaultData(r)
	data["error_code"] = code
	data["error_message"] = message
	err := templates.RenderTemplateSafe(w, "errors/error.tmpl", data)
	if err != nil {
		return err
	}
	return nil
}

const internalServerErrorResponse = `
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
